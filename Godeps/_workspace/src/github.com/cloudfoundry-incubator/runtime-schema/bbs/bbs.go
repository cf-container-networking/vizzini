package bbs

import (
	"time"

	"github.com/cloudfoundry-incubator/runtime-schema/bbs/lock_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/lrp_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/services_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/start_auction_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/stop_auction_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/task_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	steno "github.com/cloudfoundry/gosteno"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/storeadapter"
)

//Bulletin Board System/Store

type RepBBS interface {
	//services
	MaintainExecutorPresence(heartbeatInterval time.Duration, executorPresence models.ExecutorPresence) (services_bbs.Presence, <-chan bool, error)

	//task
	WatchForDesiredTask() (<-chan models.Task, chan<- bool, <-chan error)
	ClaimTask(taskGuid string, executorID string) error
	StartTask(taskGuid string, executorID string, containerHandle string) error
	CompleteTask(taskGuid string, failed bool, failureReason string, result string) error

	///lrp
	ReportActualLRPAsStarting(lrp models.ActualLRP, executorID string) error
	ReportActualLRPAsRunning(lrp models.ActualLRP, executorId string) error
	RemoveActualLRP(lrp models.ActualLRP) error
	WatchForStopLRPInstance() (<-chan models.StopLRPInstance, chan<- bool, <-chan error)
	ResolveStopLRPInstance(stopInstance models.StopLRPInstance) error
}

type ConvergerBBS interface {
	//lrp
	ConvergeLRPs()

	//start auction
	ConvergeLRPStartAuctions(kickPendingDuration time.Duration, expireClaimedDuration time.Duration)

	//stop auction
	ConvergeLRPStopAuctions(kickPendingDuration time.Duration, expireClaimedDuration time.Duration)

	//task
	ConvergeTask(timeToClaim time.Duration, converganceInterval time.Duration)

	//lock
	MaintainConvergeLock(interval time.Duration, executorID string) (disappeared <-chan bool, stop chan<- chan bool, err error)
}

type TPSBBS interface {
	//lrp
	GetActualLRPsByProcessGuid(string) ([]models.ActualLRP, error)
}

type AppManagerBBS interface {
	//lrp
	GetActualLRPsByProcessGuid(string) ([]models.ActualLRP, error)
	RequestStopLRPInstance(stopInstance models.StopLRPInstance) error
	WatchForDesiredLRPChanges() (<-chan models.DesiredLRPChange, chan<- bool, <-chan error)

	//start auction
	RequestLRPStartAuction(models.LRPStartAuction) error

	//stop auction
	RequestLRPStopAuction(models.LRPStopAuction) error

	//services
	GetAvailableFileServer() (string, error)
}

type NsyncBBS interface {
	// lrp
	DesireLRP(models.DesiredLRP) error
	RemoveDesiredLRPByProcessGuid(guid string) error
	GetAllDesiredLRPs() ([]models.DesiredLRP, error)
	ChangeDesiredLRP(change models.DesiredLRPChange) error
}

type AuctioneerBBS interface {
	//services
	GetAllExecutors() ([]models.ExecutorPresence, error)

	//start auction
	WatchForLRPStartAuction() (<-chan models.LRPStartAuction, chan<- bool, <-chan error)
	ClaimLRPStartAuction(models.LRPStartAuction) error
	ResolveLRPStartAuction(models.LRPStartAuction) error

	//stop auction
	WatchForLRPStopAuction() (<-chan models.LRPStopAuction, chan<- bool, <-chan error)
	ClaimLRPStopAuction(models.LRPStopAuction) error
	ResolveLRPStopAuction(models.LRPStopAuction) error

	//lock
	MaintainAuctioneerLock(interval time.Duration, auctioneerID string) (<-chan bool, chan<- chan bool, error)
}

type StagerBBS interface {
	//task
	WatchForCompletedTask() (<-chan models.Task, chan<- bool, <-chan error)
	DesireTask(models.Task) error
	ResolvingTask(taskGuid string) error
	ResolveTask(taskGuid string) error

	//services
	GetAvailableFileServer() (string, error)
}

type MetricsBBS interface {
	//task
	GetAllTasks() ([]models.Task, error)

	//services
	GetServiceRegistrations() (models.ServiceRegistrations, error)
}

type FileServerBBS interface {
	//services
	MaintainFileServerPresence(
		heartbeatInterval time.Duration,
		fileServerURL string,
		fileServerId string,
	) (presence services_bbs.Presence, disappeared <-chan bool, err error)
}

type RouteEmitterBBS interface {
	// lrp
	WatchForDesiredLRPChanges() (<-chan models.DesiredLRPChange, chan<- bool, <-chan error)
	WatchForActualLRPChanges() (<-chan models.ActualLRPChange, chan<- bool, <-chan error)
	GetAllDesiredLRPs() ([]models.DesiredLRP, error)
	GetRunningActualLRPs() ([]models.ActualLRP, error)
	GetDesiredLRPByProcessGuid(processGuid string) (models.DesiredLRP, error)
	GetRunningActualLRPsByProcessGuid(processGuid string) ([]models.ActualLRP, error)
}

func NewRepBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) RepBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewConvergerBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) ConvergerBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewAppManagerBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) AppManagerBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewNsyncBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) NsyncBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewAuctioneerBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) AuctioneerBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewStagerBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) StagerBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewMetricsBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) MetricsBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewFileServerBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) FileServerBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewRouteEmitterBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) RouteEmitterBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewTPSBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) TPSBBS {
	return NewBBS(store, timeProvider, logger)
}

func NewBBS(store storeadapter.StoreAdapter, timeProvider timeprovider.TimeProvider, logger *steno.Logger) *BBS {
	return &BBS{
		LockBBS:         lock_bbs.New(store),
		LRPBBS:          lrp_bbs.New(store, timeProvider, logger),
		StartAuctionBBS: start_auction_bbs.New(store, timeProvider, logger),
		StopAuctionBBS:  stop_auction_bbs.New(store, timeProvider, logger),
		ServicesBBS:     services_bbs.New(store, logger),
		TaskBBS:         task_bbs.New(store, timeProvider, logger),
	}
}

type BBS struct {
	*lock_bbs.LockBBS
	*lrp_bbs.LRPBBS
	*start_auction_bbs.StartAuctionBBS
	*stop_auction_bbs.StopAuctionBBS
	*services_bbs.ServicesBBS
	*task_bbs.TaskBBS
}
