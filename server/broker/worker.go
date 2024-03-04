package broker

// IJob : Interface for the Job to be processed
type IJob interface {
	Process() error
}

// IWorker : Interface for Worker
type IWorker interface {
	Start()
	Stop()
}

// Worker : Default Worker implementation
type Worker struct {
	WorkerPool   chan chan IJob // A pool of workers channels that are registered in the dispatcher
	JobChannel   chan IJob      // Channel through which a job is received by the worker
	Quit         chan bool      // Channel for Quit signal
	WorkerNumber int            // Worker Number
}

// Start : Start the worker and add to worker pool
func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel
			select {
			case job := <-w.JobChannel: // Worker is waiting here to receive job from JobQueue
				job.Process() // Worker is Processing the job

			case <-w.Quit:
				// Signal to stop the worker
				return
			}
		}
	}()
}

// Stop : Calling this method stops the worker
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

func newWorker(workerPool chan chan IJob, number int) IWorker {
	return &Worker{
		WorkerPool:   workerPool,
		JobChannel:   make(chan IJob),
		Quit:         make(chan bool),
		WorkerNumber: number,
	}
}

// Option sets a parameter for the Dispatcher
type Option func(d *Dispatcher)

// SetMaxWorkers sets the number of workers. Default is 10
func SetMaxWorkers(maxWorkers int) Option {
	return func(d *Dispatcher) {
		if maxWorkers > 0 {
			d.maxWorkers = maxWorkers
		}
	}
}

// SetNewWorker sets the Worker initialisation function in dispatcher
func SetNewWorker(newWorker func(chan chan IJob, int) IWorker) Option {
	return func(d *Dispatcher) {
		d.newWorker = newWorker
	}
}

// SetJobQueue sets the JobQueue in dispatcher
func SetJobQueue(jobQueue chan IJob) Option {
	return func(d *Dispatcher) {
		d.JobQueue = jobQueue
	}
}

// Dispatcher holds worker pool, job queue and manages workers and job
// To submit a job to worker pool, use code
// `dispatcher.JobQueue <- job`
type Dispatcher struct {
	name       string
	workerPool chan chan IJob // A pool of workers channels that are registered with the dispatcher
	maxWorkers int
	newWorker  func(chan chan IJob, int) IWorker
	JobQueue   chan IJob
}

func (d *Dispatcher) run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		go func(j int) {
			worker := d.newWorker(d.workerPool, j) // Initialise a new worker
			worker.Start()
		}(i) // Start the worker
	}
	go d.dispatch() // Start the dispatcher
}

func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue {
		// try to obtain a worker job channel that is available.
		// this will block until a worker is idle
		jobChannel := <-d.workerPool
		// dispatch the job to the worker job channel
		jobChannel <- job
	}
}

// NewDispatcher : returns a new dispatcher. When no options are given, it returns a dispatcher with default settings
// 10 Workers and `newWorker` initialisation
func NewDispatcher(dispatcherName string, options ...Option) *Dispatcher {
	d := &Dispatcher{
		name:       dispatcherName,
		maxWorkers: 10,
		newWorker:  newWorker,
	}

	for _, option := range options {
		option(d)
	}
	if d.JobQueue == nil {
		d.JobQueue = make(chan IJob, d.maxWorkers)
	}

	d.workerPool = make(chan chan IJob, d.maxWorkers)
	d.run()
	return d
}
