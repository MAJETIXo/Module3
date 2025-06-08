package service

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"micro-service/Internal/dto"
	"micro-service/Internal/repo"
	"micro-service/pkg/validator"
)

type Service interface {
	CreateTask(ctx *fiber.Ctx) error
	GetTask(ctx *fiber.Ctx) error
	PutTask(ctx *fiber.Ctx) error
	DeleteTask(ctx *fiber.Ctx) error
	GetTasks(ctx *fiber.Ctx) error
}

type service struct {
	repo repo.Repo
	log  *zap.SugaredLogger
}

func NewService(repo repo.Repo, logger *zap.SugaredLogger) Service {
	return &service{
		repo: repo,
		log:  logger,
	}
}

// PatchTask - обработчик запроса частичного обновления
/*func (s *service) PatchTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}
	var task = repo.TaskUpdate{
		Title:       req.Title,
		Description: req.Description,
	}
	taskID, err := s.repo.PatchTask(ctx.Context(), task)
	if err != nil {
		s.log.Error("Failed to patch task", zap.Error(err))
		return dto.IternalServerError(ctx)
	}
	response := dto.Response{
		Status: "ok",
		Data:   map[string]int{"id": taskID},
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}
*/

// CreateTask - обработчик запроса на создание задачи
func (s *service) CreateTask(ctx *fiber.Ctx) error {
	var req TaskRequest

	// Десериализация JSON-запроса
	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Вставка задачи в БД через репозиторий
	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
	}
	taskID, err := s.repo.CreateTask(ctx.Context(), task)
	if err != nil {
		s.log.Error("Failed to insert task", zap.Error(err))
		return dto.IternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": taskID},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetTask - обработчик на получение информации задачи
func (s *service) GetTask(ctx *fiber.Ctx) error {
	// Получение параметра id из URL
	id, err := ctx.ParamsInt("id")
	if err != nil {
		s.log.Error("Invalid task id(not int)", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid ID format")
	}

	// Получение задачи из БД
	task, err := s.repo.GetTask(ctx.Context(), id)
	if err != nil {
		s.log.Error("Failed to retrive task", zap.Error(err))
		return dto.NotFoundError(ctx, "Task not found")
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   task,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetTasks - обработчик на получение информации задач
func (s *service) GetTasks(ctx *fiber.Ctx) error {
	// Получение всех задач из БД
	tasks, err := s.repo.GetTasks(ctx.Context())
	if err != nil {
		s.log.Error("Failed to retrieve tasks", zap.Error(err))
		return dto.IternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   tasks, // Теперь передаём срез задач []Task
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// PutTask - обработчик запроса полного обновления
func (s *service) PutTask(ctx *fiber.Ctx) error {
	// Получение ID задачи из URL
	id, err := ctx.ParamsInt("id")
	if err != nil {
		s.log.Error("Invalid task id(not int)", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid ID format")
	}

	// Десериализация JSON-запроса
	var req TaskRequest

	if err := json.Unmarshal(ctx.Body(), &req); err != nil {
		s.log.Error("Invalid request body", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid request body")
	}

	// Валидация входных данных
	if vErr := validator.Validate(ctx.Context(), req); vErr != nil {
		return dto.BadResponseError(ctx, dto.FieldIncorrect, vErr.Error())
	}

	// Обновление задачи в БД
	task := repo.Task{
		Title:       req.Title,
		Description: req.Description,
	}
	id, err = s.repo.PutTask(ctx.Context(), id, task)
	if err != nil {
		s.log.Error("Failed to update task", zap.Error(err))
		return dto.IternalServerError(ctx)
	}

	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"task_id": id},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// DeleteTask - обработчик запроса на удаление задачи
func (s *service) DeleteTask(ctx *fiber.Ctx) error {
	// Получение ID задачи из URL
	id, err := ctx.ParamsInt("id")
	if err != nil {
		s.log.Error("Invalid task id(not int)", zap.Error(err))
		return dto.BadResponseError(ctx, dto.FieldBadFormat, "Invalid ID format")
	}

	// Удаление задачи из БД
	id, err = s.repo.DeleteTask(ctx.Context(), id)
	if err != nil {
		s.log.Error("Failed to delete task", zap.Error(err))
		return dto.NotFoundError(ctx, "Id not found")
	}
	// Формирование ответа
	response := dto.Response{
		Status: "success",
		Data:   map[string]int{"deleted_task_id": id},
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}
