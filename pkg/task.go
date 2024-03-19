package pkg

import log "github.com/sirupsen/logrus"

type StaticTaskRunner[Task any] interface {
	Generator(chan Task)
	Solver(Task)
}

func RunStatic[Task any](runner StaticTaskRunner[Task], threads int) {
	log.Debugf("RunStatic threads: %v", threads)
	tasks := make(chan Task, threads*3)
	for i := 0; i < threads; i++ {
		go func() {
			for task := range tasks {
				runner.Solver(task)
			}
		}()
	}
	runner.Generator(tasks)
	close(tasks)
}
