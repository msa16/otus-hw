package hw06pipelineexecution

import (
	"sync/atomic"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	stopFlag := atomic.Bool{}
	for _, stage := range stages {
		// промежуточный канал чтобы мы могли его закрыть, и stage завершил внутренний цикл получения и обработки
		proxyIn := make(Bi)
		out := stage(proxyIn)
		go func(stageIn In) {
			defer func() {
				// закрываем промежуточный канал чтобы наш stage завершил цикл потому что ему нечего читать
				close(proxyIn)
				// вычитать данные от предыдущего stage чтобы он не остановился на записи к нам
				//nolint:revive
				for range stageIn {
				}
			}()
			for {
				// сигнал done дойдёт до всех горутин, но если готовы события и чтения, и остановки, то порядок
				// их выбора не гарантирован, поэтому прерывание слушают все, выставляют флаг, и все также смотрят
				// на флаг, который могла выставить другая горутина
				if stopFlag.Load() {
					return
				}
				select {
				case <-done:
					stopFlag.Store(true)
					return
				case v, ok := <-stageIn:
					// если не смогли прочитать, или прочитали, но кто-то взвёл флаг остановки, то не будем
					// передавать данные в stage, а сразу уходим
					if !ok || stopFlag.Load() {
						return
					}
					// что быстрее - или запишем данные в stage, или придёт сигнал остановки
					select {
					case <-done:
						stopFlag.Store(true)
						return
					case proxyIn <- v:
					}
				}
			}
		}(in)
		in = out
	}
	return in
}
