package handlers

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wailorman/fftb/pkg/distributed/remote/pb"
)

type dealerHelper struct {
	dealer        pb.Dealer
	ctx           context.Context
	authorization string
	logger        *logrus.Entry
	task          *pb.Task
}

func newDealerHelper(handler interface{}) *dealerHelper {
	h := &dealerHelper{}

	switch th := handler.(type) {
	case *ConvertTaskHandler:
		h.dealer = th.dealer
		h.ctx = th.ctx
		h.authorization = th.authorization
		h.logger = th.logger
		h.task = th.task

	case *MediaMetaHandler:
		h.dealer = th.dealer
		h.ctx = th.ctx
		h.authorization = th.authorization
		h.logger = th.logger
		h.task = th.task

	default:
		panic("Unknown handler type")
	}

	return h
}

func (h *dealerHelper) exit(err error) {
	if errors.Is(h.ctx.Err(), context.Canceled) || errors.Is(err, context.Canceled) {
		h.quit()
		return
	}

	if err != nil {
		h.fail(err)
		return
	}

	h.finish()
}

func (h *dealerHelper) finish() {
	_, err := h.dealer.FinishTask(context.TODO(), &pb.FinishTaskRequest{
		Authorization: h.authorization,
		TaskId:        h.task.Id,
	})

	if err != nil {
		h.logger.WithError(err).Warn("Failed to finish")
	}
}

func (h *dealerHelper) quit() {
	_, err := h.dealer.QuitTask(context.TODO(), &pb.QuitTaskRequest{
		Authorization: h.authorization,
		TaskId:        h.task.Id,
	})

	if err != nil {
		h.logger.WithError(err).Warn("Failed to quit")
	}
}

func (h *dealerHelper) fail(failure error) {
	_, err := h.dealer.FailTask(h.ctx, failTaskRequest(h.authorization, h.task.Id, failure))

	if err != nil {
		h.logger.WithError(err).Warn("Failed to report failure")
	}
}

func failTaskRequest(authorization string, taskId string, err error) *pb.FailTaskRequest {
	return &pb.FailTaskRequest{
		Authorization: authorization,
		TaskId:        taskId,
		Failures:      []string{err.Error()},
	}
}
