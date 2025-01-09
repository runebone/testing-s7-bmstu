package v1_test

import (
	log "aggregator/internal/adapter/logger"
	"aggregator/internal/dto"
	"aggregator/mocks"
	"context"
	"errors"
	"testing"
	"time"

	v1 "aggregator/internal/usecase/v1"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func TestGetStats(t *testing.T) {
	runner.Run(t, "TestGetStats", func(pt provider.T) {
		fromTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
		toTime := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)

		tests := []struct {
			name      string
			from      time.Time
			to        time.Time
			mockSetup func(mockUserSvc *mocks.UserService, mockTodoSvc *mocks.TodoService, from, to time.Time)
			wantErr   bool
			err       error
		}{
			{
				name: "positive",
				from: fromTime,
				to:   toTime,
				mockSetup: func(mockUserSvc *mocks.UserService, mockTodoSvc *mocks.TodoService, from, to time.Time) {
					userDTOs := make([]dto.User, 3)

					userDTOs[0] = dto.User{
						ID:       uuid.New(),
						Email:    "user@zero.com",
						Username: "UserZero",
					}
					userDTOs[1] = dto.User{
						ID:       uuid.New(),
						Email:    "user@one.com",
						Username: "UserOne",
					}
					userDTOs[2] = dto.User{
						ID:       uuid.New(),
						Email:    "user@two.com",
						Username: "UserTwo",
					}

					mockUserSvc.On("GetNewUsers", context.Background(), from, to).Return(userDTOs, nil)

					cardDTOs := make([]dto.Card, 3)

					cardDTOs[0] = dto.Card{
						ID:       uuid.New(),
						UserID:   uuid.New(),
						ColumnID: uuid.New(),
						Title:    "TitleZero",
					}
					cardDTOs[1] = dto.Card{
						ID:       uuid.New(),
						UserID:   userDTOs[1].ID,
						ColumnID: uuid.New(),
						Title:    "TitleOne",
					}
					cardDTOs[2] = dto.Card{
						ID:       uuid.New(),
						UserID:   userDTOs[2].ID,
						ColumnID: uuid.New(),
						Title:    "TitleTwo",
					}

					mockTodoSvc.On("GetNewCards", context.Background(), from, to).Return(cardDTOs, nil)
				},
				wantErr: false,
			},
			{
				name:      "negative",
				from:      toTime,
				to:        fromTime,
				mockSetup: func(mockUserSvc *mocks.UserService, mockTodoSvc *mocks.TodoService, from, to time.Time) {},
				wantErr:   true,
				err:       v1.ErrInvalidTimeRange,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockUserSvc, mockTodoSvc, tt.from, tt.to)

					pt.WithNewStep("Call GetStats", func(sCtx provider.StepCtx) {
						_, err := uc.GetStats(context.Background(), tt.from, tt.to)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockUserSvc.AssertExpectations(t)
						mockTodoSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestRegister(t *testing.T) {
	runner.Run(t, "TestRegister", func(pt provider.T) {
		tests := []struct {
			name      string
			username  string
			email     string
			password  string
			mockSetup func(mockAuthSvc *mocks.AuthService, username, email, password string)
			wantErr   bool
			err       error
		}{
			{
				name:     "positive",
				username: "PositiveUsername",
				email:    "positive@email.com",
				password: "P0s1t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, username, email, password string) {
					tokens := dto.Tokens{
						AccessToken:  "AccessToken",
						RefreshToken: "RefreshToken",
					}

					mockAuthSvc.On("Register", context.Background(), username, email, password).Return(&tokens, nil)
				},
				wantErr: false,
			},
			{
				name:     "negative",
				username: "NegativeUsername",
				email:    "negative@email.com",
				password: "N3g@t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, username, email, password string) {
					mockAuthSvc.On("Register", context.Background(), username, email, password).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrRegister,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.username, tt.email, tt.password)

					pt.WithNewStep("Call Register", func(sCtx provider.StepCtx) {
						_, err := uc.Register(context.Background(), tt.username, tt.email, tt.password)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestLogin(t *testing.T) {
	runner.Run(t, "TestLogin", func(pt provider.T) {
		tests := []struct {
			name      string
			email     string
			password  string
			mockSetup func(mockAuthSvc *mocks.AuthService, email, password string)
			wantErr   bool
			err       error
		}{
			{
				name:     "positive",
				email:    "positive@email.com",
				password: "P0s1t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, email, password string) {
					tokens := dto.Tokens{
						AccessToken:  "AccessToken",
						RefreshToken: "RefreshToken",
					}

					mockAuthSvc.On("Login", context.Background(), email, password).Return(&tokens, nil)
				},
				wantErr: false,
			},
			{
				name:     "negative",
				email:    "negative@email.com",
				password: "N3g@t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, email, password string) {
					mockAuthSvc.On("Login", context.Background(), email, password).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrLogin,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.email, tt.password)

					pt.WithNewStep("Call Login", func(sCtx provider.StepCtx) {
						_, err := uc.Login(context.Background(), tt.email, tt.password)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestRefresh(t *testing.T) {
	runner.Run(t, "TestRefresh", func(pt provider.T) {
		tests := []struct {
			name         string
			refreshToken string
			mockSetup    func(mockAuthSvc *mocks.AuthService, refreshToken string)
			wantErr      bool
			err          error
		}{
			{
				name:         "positive",
				refreshToken: "PositiveToken",
				mockSetup: func(mockAuthSvc *mocks.AuthService, refreshToken string) {
					resp := dto.RefreshResponse{
						AccessToken: "AccessToken",
					}

					mockAuthSvc.On("Refresh", context.Background(), refreshToken).Return(&resp, nil)
				},
				wantErr: false,
			},
			{
				name:         "negative",
				refreshToken: "NegativeToken",
				mockSetup: func(mockAuthSvc *mocks.AuthService, refreshToken string) {
					mockAuthSvc.On("Refresh", context.Background(), refreshToken).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrRefresh,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.refreshToken)

					pt.WithNewStep("Call Refresh", func(sCtx provider.StepCtx) {
						_, err := uc.Refresh(context.Background(), tt.refreshToken)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestValidate(t *testing.T) {
	runner.Run(t, "TestValidate", func(pt provider.T) {
	})
}

func TestLogout(t *testing.T) {
	runner.Run(t, "TestLogout", func(pt provider.T) {
	})
}

func TestGetBoards(t *testing.T) {
	runner.Run(t, "TestGetBoards", func(pt provider.T) {
	})
}

func TestGetColumns(t *testing.T) {
	runner.Run(t, "TestGetColumns", func(pt provider.T) {
	})
}

func TestGetCards(t *testing.T) {
	runner.Run(t, "TestGetCards", func(pt provider.T) {
	})
}

func TestGetCard(t *testing.T) {
	runner.Run(t, "TestGetCard", func(pt provider.T) {
	})
}

func TestCreateBoard(t *testing.T) {
	runner.Run(t, "TestCreateBoard", func(pt provider.T) {
	})
}

func TestCreateColumn(t *testing.T) {
	runner.Run(t, "TestCreateColumn", func(pt provider.T) {
	})
}

func TestCreateCard(t *testing.T) {
	runner.Run(t, "TestCreateCard", func(pt provider.T) {
	})
}

func TestUpdateBoard(t *testing.T) {
	runner.Run(t, "TestUpdateBoard", func(pt provider.T) {
	})
}

func TestUpdateColumn(t *testing.T) {
	runner.Run(t, "TestUpdateColumn", func(pt provider.T) {
	})
}

func TestUpdateCard(t *testing.T) {
	runner.Run(t, "TestUpdateCard", func(pt provider.T) {
	})
}

func TestDeleteBoard(t *testing.T) {
	runner.Run(t, "TestDeleteBoard", func(pt provider.T) {
	})
}

func TestDeleteColumn(t *testing.T) {
	runner.Run(t, "TestDeleteColumn", func(pt provider.T) {
	})
}

func TestDeleteCard(t *testing.T) {
	runner.Run(t, "TestDeleteCard", func(pt provider.T) {
		tests := []struct {
			name      string
			mockSetup func()
			wantErr   bool
			err       error
		}{
			{
				name:      "positive",
				mockSetup: func() {},
				wantErr:   false,
			},
			{
				name:      "negative",
				mockSetup: func() {},
				wantErr:   true,
				// err:       ,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					// mockUserSvc := new(mocks.UserService)
					// mockAuthSvc := new(mocks.AuthService)
					// mockTodoSvc := new(mocks.TodoService)
					// logger := log.NewEmptyLogger()

					// uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					// tt.mockSetup()

					pt.WithNewStep("Call TODO", func(sCtx provider.StepCtx) {
						// err := userUC.CreateUser(context.Background(), tt.user)

						if tt.wantErr {
							// sCtx.Assert().Error(err, "Expected error")
							// sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							// sCtx.Assert().NoError(err, "Expected no error")
						}

						// mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

// package v1_test
//
// import (
// 	"aggregator/internal/dto"
// 	"aggregator/internal/entity"
// 	"aggregator/internal/usecase"
// 	v1 "aggregator/internal/usecase/v1"
// 	"aggregator/mocks"
// 	"context"
// 	"errors"
// 	"sort"
// 	"testing"
// 	"time"
//
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )
//
// type testSetup struct {
// 	ctx         context.Context
// 	mockUserSvc *mocks.UserService
// 	mockAuthSvc *mocks.AuthService
// 	mockTodoSvc *mocks.TodoService
// 	uc          usecase.AggregatorUseCase
// }
//
// func setup() *testSetup {
// 	ctx := context.TODO()
//
// 	mockUserSvc := new(mocks.UserService)
// 	mockAuthSvc := new(mocks.AuthService)
// 	mockTodoSvc := new(mocks.TodoService)
//
// 	aggregatorUseCase := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc)
//
// 	return &testSetup{
// 		ctx:         ctx,
// 		mockUserSvc: mockUserSvc,
// 		mockAuthSvc: mockAuthSvc,
// 		mockTodoSvc: mockTodoSvc,
// 		uc:          aggregatorUseCase,
// 	}
// }
//
// func Date(yyyy, mm, dd int) time.Time {
// 	return time.Date(yyyy, time.Month(mm), dd, 0, 0, 0, 0, time.UTC)
// }
//
// func TestGetStats(t *testing.T) {
// 	ts := setup()
//
// 	mockUserFnOk := func(from, to time.Time, users []dto.User) {
// 		ts.mockUserSvc.On("GetNewUsers", ts.ctx, from, to).Return(users, nil)
// 	}
//
// 	mockUserFnErr := func(from, to time.Time, users []dto.User) {
// 		ts.mockUserSvc.On("GetNewUsers", ts.ctx, from, to).Return(nil, errors.New(""))
// 	}
//
// 	mockTodoFnOk := func(from, to time.Time, cards []dto.Card) {
// 		ts.mockTodoSvc.On("GetNewCards", ts.ctx, from, to).Return(cards, nil)
// 	}
//
// 	mockTodoFnErr := func(from, to time.Time, cards []dto.Card) {
// 		ts.mockTodoSvc.On("GetNewCards", ts.ctx, from, to).Return(nil, errors.New(""))
// 	}
//
// 	dates := []time.Time{
// 		Date(2000, 1, 1),
// 		Date(2000, 1, 2),
// 		Date(2000, 1, 3),
// 		Date(2000, 1, 4),
// 	}
//
// 	nu := 4
// 	uuids := make([]uuid.UUID, nu)
// 	for i := 0; i < nu; i++ {
// 		uuids[i] = uuid.New()
// 	}
//
// 	users := make([]dto.User, nu)
// 	for i := 0; i < nu; i++ {
// 		users[i] = dto.User{
// 			ID: uuids[i],
// 		}
// 	}
//
// 	users[0].CreatedAt = dates[0]
// 	users[1].CreatedAt = dates[1]
// 	users[2].CreatedAt = dates[1]
// 	users[3].CreatedAt = dates[2]
//
// 	nc := 13
// 	cards := make([]dto.Card, nc)
// 	for i := 0; i < nc; i++ {
// 		cards[i].ID = uuid.New()
// 	}
//
// 	// User0:
// 	// - Day0:
// 	//   - Card0
// 	cards[0].UserID = uuids[0]
// 	cards[0].CreatedAt = dates[0]
//
// 	// - Day1:
// 	//   - Card1
// 	cards[1].UserID = uuids[0]
// 	cards[1].CreatedAt = dates[1]
//
// 	//   - Card2
// 	cards[2].UserID = uuids[0]
// 	cards[2].CreatedAt = dates[1]
//
// 	// - Day2:
// 	//   - Card6
// 	cards[6].UserID = uuids[0]
// 	cards[6].CreatedAt = dates[2]
//
// 	//   - Card7
// 	cards[7].UserID = uuids[0]
// 	cards[7].CreatedAt = dates[2]
//
// 	//   - Card8
// 	cards[8].UserID = uuids[0]
// 	cards[8].CreatedAt = dates[2]
//
// 	// - Day3:
// 	//   - Card11
// 	cards[11].UserID = uuids[0]
// 	cards[11].CreatedAt = dates[3]
//
// 	//   - Card12
// 	cards[12].UserID = uuids[0]
// 	cards[12].CreatedAt = dates[3]
//
// 	// User1:
// 	// - Day1:
// 	//   - Card3
// 	cards[3].UserID = uuids[1]
// 	cards[3].CreatedAt = dates[1]
//
// 	//   - Card4
// 	cards[4].UserID = uuids[1]
// 	cards[4].CreatedAt = dates[1]
//
// 	//   - Card5
// 	cards[5].UserID = uuids[1]
// 	cards[5].CreatedAt = dates[1]
//
// 	// User2:
// 	// - Day1:
// 	//   - Card9
// 	cards[9].UserID = uuids[2]
// 	cards[9].CreatedAt = dates[1]
//
// 	//   - Card10
// 	cards[10].UserID = uuids[2]
// 	cards[10].CreatedAt = dates[1]
//
// 	// User3:
// 	// - Day2:
//
// 	tests := []struct {
// 		name       string
// 		from, to   time.Time
// 		users      []dto.User
// 		cards      []dto.Card
// 		stats      []entity.NewUsersAndCardsStats
// 		mockUserFn func(from, to time.Time, users []dto.User)
// 		mockTodoFn func(from, to time.Time, cards []dto.Card)
// 		wantErr    bool
// 		errMsg     string
// 	}{
// 		{
// 			name: "success, first",
// 			from: dates[0],
// 			to:   dates[0],
// 			users: []dto.User{
// 				users[0],
// 			},
// 			cards: []dto.Card{
// 				cards[0],
// 			},
// 			stats: []entity.NewUsersAndCardsStats{
// 				{
// 					Date: dates[0],
// 					Users: dto.ToUserEntities([]dto.User{
// 						users[0],
// 					}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[0],
// 					}),
// 					NumCardsByNewUsers: 1,
// 				},
// 			},
// 			mockUserFn: mockUserFnOk,
// 			mockTodoFn: mockTodoFnOk,
// 			wantErr:    false,
// 		},
// 		{
// 			name:  "success, all",
// 			from:  dates[0],
// 			to:    dates[len(dates)-1],
// 			users: users,
// 			cards: cards,
// 			stats: []entity.NewUsersAndCardsStats{
// 				{
// 					Date: dates[0],
// 					Users: dto.ToUserEntities([]dto.User{
// 						users[0],
// 					}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[0],
// 					}),
// 					NumCardsByNewUsers: 1,
// 				},
// 				{
// 					Date: dates[1],
// 					Users: dto.ToUserEntities([]dto.User{
// 						users[1],
// 						users[2],
// 					}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[1],
// 						cards[2],
// 						cards[3],
// 						cards[4],
// 						cards[5],
// 						cards[9],
// 						cards[10],
// 					}),
// 					NumCardsByNewUsers: 5,
// 				},
// 				{
// 					Date: dates[2],
// 					Users: dto.ToUserEntities([]dto.User{
// 						users[3],
// 					}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[6],
// 						cards[7],
// 						cards[8],
// 					}),
// 					NumCardsByNewUsers: 0,
// 				},
// 				{
// 					Date:  dates[3],
// 					Users: dto.ToUserEntities([]dto.User{}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[11],
// 						cards[12],
// 					}),
// 					NumCardsByNewUsers: 0,
// 				},
// 			},
// 			mockUserFn: mockUserFnOk,
// 			mockTodoFn: mockTodoFnOk,
// 			wantErr:    false,
// 		},
// 		{
// 			name: "success, some",
// 			from: dates[1],
// 			to:   dates[len(dates)-2],
// 			users: []dto.User{
// 				users[1],
// 				users[2],
// 				users[3],
// 			},
// 			cards: []dto.Card{
// 				cards[1],
// 				cards[2],
// 				cards[3],
// 				cards[4],
// 				cards[5],
// 				cards[6],
// 				cards[7],
// 				cards[8],
// 				cards[9],
// 				cards[10],
// 			},
// 			stats: []entity.NewUsersAndCardsStats{
// 				{
// 					Date: dates[1],
// 					Users: dto.ToUserEntities([]dto.User{
// 						users[1],
// 						users[2],
// 					}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[1],
// 						cards[2],
// 						cards[3],
// 						cards[4],
// 						cards[5],
// 						cards[9],
// 						cards[10],
// 					}),
// 					NumCardsByNewUsers: 5,
// 				},
// 				{
// 					Date: dates[2],
// 					Users: dto.ToUserEntities([]dto.User{
// 						users[3],
// 					}),
// 					Cards: dto.ToCardEntities([]dto.Card{
// 						cards[6],
// 						cards[7],
// 						cards[8],
// 					}),
// 					NumCardsByNewUsers: 0,
// 				},
// 			},
// 			mockUserFn: mockUserFnOk,
// 			mockTodoFn: mockTodoFnOk,
// 			wantErr:    false,
// 		},
// 		{
// 			name:       "fail - <<from>> is greater than <<to>>",
// 			from:       dates[1],
// 			to:         dates[0],
// 			users:      users,
// 			cards:      cards,
// 			stats:      []entity.NewUsersAndCardsStats{},
// 			mockUserFn: func(from, to time.Time, users []dto.User) {},
// 			mockTodoFn: func(from, to time.Time, users []dto.Card) {},
// 			wantErr:    true,
// 			errMsg:     v1.ErrInvalidTimeRange.Error(),
// 		},
// 		{
// 			name:       "failed to get new users",
// 			from:       dates[0],
// 			to:         dates[1],
// 			users:      users,
// 			cards:      cards,
// 			stats:      []entity.NewUsersAndCardsStats{},
// 			mockUserFn: mockUserFnErr,
// 			mockTodoFn: mockTodoFnOk,
// 			wantErr:    true,
// 			errMsg:     v1.ErrGetNewUsers.Error(),
// 		},
// 		{
// 			name:       "failed to get new cards",
// 			from:       dates[0],
// 			to:         dates[1],
// 			users:      users,
// 			cards:      cards,
// 			stats:      []entity.NewUsersAndCardsStats{},
// 			mockUserFn: mockUserFnOk,
// 			mockTodoFn: mockTodoFnErr,
// 			wantErr:    true,
// 			errMsg:     v1.ErrGetNewCards.Error(),
// 		},
// 	}
//
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.mockUserFn(tt.from, tt.to, tt.users)
// 			tt.mockTodoFn(tt.from, tt.to, tt.cards)
//
// 			stats, err := ts.uc.GetStats(ts.ctx, tt.from, tt.to)
//
// 			csExpected := ToComparableStats(tt.stats)
// 			csActual := ToComparableStats(stats)
//
// 			if tt.wantErr {
// 				assert.NotNil(t, err)
// 				assert.EqualError(t, err, tt.errMsg)
// 			} else {
// 				assert.Nil(t, err)
// 				// assert.Equal(t, tt.stats, stats)
// 				assert.Equal(t, csExpected, csActual)
// 				ts.mockUserSvc.AssertCalled(t, "GetNewUsers", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
// 				ts.mockTodoSvc.AssertCalled(t, "GetNewCards", ts.ctx, mock.Anything, mock.Anything, mock.Anything)
// 			}
//
// 			t.Cleanup(func() {
// 				ts.mockUserSvc.ExpectedCalls = nil
// 				ts.mockUserSvc.Calls = nil
//
// 				ts.mockTodoSvc.ExpectedCalls = nil
// 				ts.mockTodoSvc.Calls = nil
// 			})
// 		})
// 	}
// }
//
// type ComparableStats struct {
// 	Date               time.Time
// 	NumUsers           int
// 	NumCards           int
// 	NumCardsByNewUsers int
// }
//
// func ToComparableStats(stats []entity.NewUsersAndCardsStats) []ComparableStats {
// 	cs := make([]ComparableStats, len(stats))
//
// 	for i, stat := range stats {
// 		cs[i] = ComparableStats{
// 			Date:               stat.Date,
// 			NumUsers:           len(stat.Users),
// 			NumCards:           len(stat.Cards),
// 			NumCardsByNewUsers: stat.NumCardsByNewUsers,
// 		}
// 	}
//
// 	sort.Slice(cs, func(i, j int) bool {
// 		return cs[i].Date.Unix() < cs[j].Date.Unix()
// 	})
//
// 	return c
