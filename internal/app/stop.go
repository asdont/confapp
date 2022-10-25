package app

import (
	"os"
	"os/signal"
	"syscall"
)

func SysStop(chStop chan<- struct{}) {
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan,
		syscall.SIGHUP,  // Отключение процесса от родительского процесса.
		syscall.SIGINT,  // Прерывание с клавиатуры(ctrl+c).
		syscall.SIGQUIT, // Выход с клавиатуры(Ctrl+\, Ctrl+4, SysRq).
		syscall.SIGTERM, // Завершение работы.
	)

	<-sigChan

	chStop <- struct{}{}
}
