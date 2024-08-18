package sample

import (
	"errors"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apihandler"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/customerrors"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/domain/sample"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/repository"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/validate"
	"github.com/gofiber/fiber/v2"
)

type SampleAPIHandler struct {
	sampleRepository repository.SampleRepository
}

func NewSampleAPIHandler(sampleRepository repository.SampleRepository) *SampleAPIHandler {
	return &SampleAPIHandler{sampleRepository: sampleRepository}
}

func (h *SampleAPIHandler) Pattern() string {
	return "/samples"
}

func (h *SampleAPIHandler) Handlers() []*apihandler.APIHandler {
	return []*apihandler.APIHandler{
		// 가독성 좋게 함수형태로 변경
		{
			Pattern: "",
			Method:  fiber.MethodGet,
			Handler: h.listSamples(),
		},
		{
			Pattern: "/:id",
			Method:  fiber.MethodGet,
			Handler: h.getSample(),
		},
		{
			Pattern: "",
			Method:  fiber.MethodPost,
			Handler: h.createSample(),
		},
		{
			Pattern: "/:id",
			Method:  fiber.MethodPut,
			Handler: h.updateSample(),
		},
		{
			Pattern: "/:id",
			Method:  fiber.MethodDelete,
			Handler: h.deleteSample(),
		},
	}
}

func (h *SampleAPIHandler) listSamples() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		samples, err := h.sampleRepository.FindAllSamples()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&ListSamplesResponse{Samples: samples})
	}
}

func (h *SampleAPIHandler) getSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")
		sample, err := h.sampleRepository.FindOneSample(sampleID)
		if err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&GetSampleResponse{Sample: sample})
	}
}

func (h *SampleAPIHandler) deleteSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")
		if err := h.sampleRepository.DeleteSample(sampleID); err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&DeleteSampleResponse{})
	}
}

func (h *SampleAPIHandler) createSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := new(CreateSampleRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		// 유효성 검사
		if err := validate.Validator.Struct(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		sample, err := h.sampleRepository.CreateSample(&sample.Sample{
			Name:  request.Name,
			Email: request.Email,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&CreateSampleResponse{Sample: sample})
	}
}

func (h *SampleAPIHandler) updateSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")

		request := new(UpdateSampleRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		// 유효성 검사
		if err := validate.Validator.Struct(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		sample, err := h.sampleRepository.UpdateSample(
			&sample.Sample{
				ID:    sampleID,
				Name:  request.Name,
				Email: request.Email,
			},
		)
		if err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&UpdateSampleResponse{
			Sample: sample,
		})
	}
}
