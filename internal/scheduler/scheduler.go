package scheduler

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/superbet-group/logging.clients"

	"github.com/superbet-group/offer.fastly.sse-load-test/internal/models"
)

type Scheduler struct {
	log         logging.Logger
	offerClient OfferClient
	sseClient   SSEClient

	noOfMatches         int
	connectionsPerMatch int

	rescheduleMatchesChan chan struct{}

	matchWorkers map[int64][]matchWorker

	wg *sync.WaitGroup
}

func NewScheduler(log logging.Logger, offerClient OfferClient, sseClient SSEClient, noOfMatches, connectionsPerMatch int) *Scheduler {
	return &Scheduler{
		log:                   log,
		offerClient:           offerClient,
		sseClient:             sseClient,
		noOfMatches:           noOfMatches,
		connectionsPerMatch:   connectionsPerMatch,
		rescheduleMatchesChan: make(chan struct{}, 1),
		matchWorkers:          map[int64][]matchWorker{},
		wg:                    new(sync.WaitGroup),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	rand.Seed(time.Now().Unix())

	s.rescheduleMatchesChan <- struct{}{}

	for {
		select {
		case <-ctx.Done():
			return

		case <-s.rescheduleMatchesChan:
			s.stopOldWorkers()
			offer, err := s.offerClient.GetOffer()
			if err != nil {
				s.log.Error("error while fetching offer", logging.Data{"error": err})
				continue
			}

			offerLen := len(offer.Data)
			ids := make([]int64, s.noOfMatches)
			indexSet := map[int64]struct{}{}
			for i := 0; i < s.noOfMatches; i++ {
				index := rand.Intn(offerLen)
				for {
					_, ok := indexSet[int64(index)]
					if ok {
						index = rand.Intn(offerLen)
						continue
					}
					break
				}
				indexSet[int64(index)] = struct{}{}
				ids[i] = offer.Data[index].ID
			}

			s.startNewWorkers(ctx, ids)
		}
	}
}

func (s *Scheduler) stopOldWorkers() {
	for _, mws := range s.matchWorkers {
		for _, mw := range mws {
			mw.Stop()
		}
	}
	s.wg.Wait()

	s.matchWorkers = map[int64][]matchWorker{}
}

func (s *Scheduler) startNewWorkers(ctx context.Context, ids []int64) {
	for _, id := range ids {
		for i := 0; i < s.connectionsPerMatch; i++ {
			worker := matchWorker{
				s.log, id, s.sseClient, make(chan struct{}), s.wg,
			}

			s.matchWorkers[id] = append(s.matchWorkers[id], worker)
			go worker.Start(ctx)
			s.wg.Add(1)
		}
	}
}

type OfferClient interface {
	GetOffer() (*models.Offer, error)
}

type SSEClient interface {
	StreamMatch(id int64) error
	StreamSport(id int64) error
	StreamTournaments(ids ...int64) error
	StreamAll() error
	StreamLive() error
	StreamPrematch() error
	StreamLiveCount() error
	StreamStructure() error
}

type matchWorker struct {
	log       logging.Logger
	matchID   int64
	sseClient SSEClient
	doneChan  chan struct{}
	wg        *sync.WaitGroup
}

func (mw *matchWorker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-mw.doneChan:
			return
		default:
			err := mw.sseClient.StreamMatch(mw.matchID)
			if err != nil {
				//mw.log.Error("match streaming error", logging.Data{"error": err})
			}

			time.Sleep(1 * time.Second)
		}
	}

	mw.wg.Done()
}

func (mw *matchWorker) Stop() {
	mw.doneChan <- struct{}{}
}
