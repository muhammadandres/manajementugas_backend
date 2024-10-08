package repository

import (
	"manajemen_tugas_master/model/domain"

	"gorm.io/gorm"
)

type TaskAndOwnerRepository interface {
	Create(user *domain.User, task *domain.Task, board *domain.Board) (*domain.Task, *domain.Owner, error)
	FindById(id uint) (*domain.TaskWithInvitation, error)
	FindAll() ([]*domain.TaskWithInvitation, error)
	FindAllOwners() ([]*domain.Task, error)
	FindAllManagers() ([]*domain.Task, error)
	FindAllEmployees() ([]*domain.Task, error)
	FindAllPlanningFiles() ([]*domain.Task, error)
	FindAllProjectFiles() ([]*domain.Task, error)
	GetNameEmailsDescription(taskID uint64) (ownerEmail string, managerEmails []string, employeeEmails []string, nametask string, description string, err error)
	Update(task *domain.Task, manager *domain.Manager, employee *domain.Employee, PlanningDescriptionFile *domain.PlanningDescriptionFile, planningFile *domain.PlanningFile, projectFile *domain.ProjectFile) (*domain.Task, *domain.Manager, *domain.Employee, *domain.PlanningDescriptionFile, *domain.PlanningFile, *domain.ProjectFile, *domain.Invitation, *domain.Invitation, error)
	UpdateValidationOwner(taskID uint, userID uint) error
	UpdateValidationManager(taskID uint, userID uint) error
	UpdateValidationEmployee(taskID uint, userID uint) error
	DeleteManager(taskId uint, managerId uint) (*gorm.DB, int64, int64, int64, error)
	DeleteEmployee(taskId uint, employeeId uint) (*gorm.DB, int64, error)
	DeletePlanningDescriptionFile(fileId uint) (*gorm.DB, string, error)
	DeletePlanningFile(fileId uint) (*gorm.DB, string, error)
	DeleteProjectFile(fileId uint) (*gorm.DB, string, error)
	Delete(taskID uint) (*gorm.DB, int64, int64, int64, int64, int64, int64, error)

	CreateInvitation(invitation *domain.Invitation) (*domain.Invitation, error)
	FindInvitationByID(id uint64) (*domain.Invitation, error)
	UpdateInvitation(invitation *domain.Invitation) error
	GetAllInvitations() ([]domain.Invitation, error)
}
