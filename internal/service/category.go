package service

import (
	"context"
	"io"

	"github.com/DiegoSenaa/go-gprc/internal/database"
	"github.com/DiegoSenaa/go-gprc/internal/pb"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, blank *pb.Blank) (*pb.ListCategoryResponse, error) {

	rows, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, err
	}

	list := make([]*pb.Category, len(rows))
	for _, category := range rows {
		value := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		list = append(list, value)
	}

	return &pb.ListCategoryResponse{
		Categories: list,
	}, nil
}

func (c *CategoryService) GetCategories(ctx context.Context, in *pb.GetCategoryRequest) (*pb.Category, error) {
	row, err := c.CategoryDB.FindByCategoryID(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{Id: row.ID, Name: row.Name, Description: row.Description}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.ListCategoryResponse{}

	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{Id: categoryResult.ID, Name: categoryResult.Name, Description: category.Description})
	}
}

func (c *CategoryService) CreateCategoryStreamBiDirectional(stream pb.CategoryService_CreateCategoryStreamBiDirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})
		if err != nil {
			return err
		}
	}
}
