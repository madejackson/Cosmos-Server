package proxy

import (
	"github.com/madejackson/cosmos-server/src/utils"
	"bufio"
	"net"
	"net/http"
	"time"
	"fmt"
	"errors"
)

type SmartResponseWriterWrapper struct {
	http.ResponseWriter
	ClientID string
	Status   int
	Bytes    int64
	ThrottleNext int
	TimeStarted time.Time
	TimeEnded time.Time
	RequestCost int
	Method string
	shield *smartShieldState
	policy utils.SmartShieldPolicy
	isOver bool
	hasBeenInterrupted bool
	isPrivileged bool
	shieldID string
}

func (w *SmartResponseWriterWrapper) IsOver() bool {
	return w.isOver
}

func (w *SmartResponseWriterWrapper) IsOld() bool {
	if !w.IsOver() {
		return false
	}
	oneHourAgo := time.Now().Add(-time.Hour)
	if w.TimeEnded.Before(oneHourAgo) {
		return true
	}
	return false
}

func (w *SmartResponseWriterWrapper) WriteHeader(status int) {
	w.Status = status
	w.RequestCost = 1
	if w.Method != "GET" {
		w.RequestCost = 5
	}
	if w.Status >= 400 {
		w.RequestCost *= 30
	}
	if !w.IsOver() {
		w.ResponseWriter.WriteHeader(status)
	}
}

func (w *SmartResponseWriterWrapper) Write(p []byte) (int, error) {
	userConsumed := shield.GetUserUsedBudgets(w.shieldID, w.ClientID)
	if !w.isPrivileged && !shield.isAllowedToReqest(w.shieldID, w.policy, userConsumed) {
		utils.Log(fmt.Sprintf("SmartShield: %s has been blocked due to abuse", w.ClientID))
		w.isOver = true
		w.TimeEnded = time.Now()
		w.hasBeenInterrupted = true
		w.ResponseWriter.WriteHeader(http.StatusServiceUnavailable)
		w.ResponseWriter.(http.Flusher).Flush()
		return 0, errors.New("Pending request cancelled due to SmartShield")
	}
	thro := 0

	if !w.isPrivileged {
		shield.computeThrottle(w.policy, userConsumed)
	}

	// initial throttle
	if w.ThrottleNext > 0 {
		time.Sleep(time.Duration(w.ThrottleNext) * time.Millisecond)
	}
	w.ThrottleNext = 0

	// ongoing throttle
	if thro > 0 {
		time.Sleep(time.Duration(thro) * time.Millisecond)
	}
	
	n, err := w.ResponseWriter.Write(p)

	if err != nil {
		w.isOver = true
		w.TimeEnded = time.Now()
		w.hasBeenInterrupted = true
	}

	w.Bytes += int64(n)
	return n, err
}

func (w *SmartResponseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hijacker.Hijack()
}

func (w *SmartResponseWriterWrapper) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		flusher.Flush()
	}
}

