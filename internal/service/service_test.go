package service_test

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/mock/gomock"

	"github.com/Nikola-Milovic/vyking-interview/internal/clients/mock"
	"github.com/Nikola-Milovic/vyking-interview/internal/domain"
	"github.com/Nikola-Milovic/vyking-interview/internal/service"
	"github.com/Nikola-Milovic/vyking-interview/internal/store"
	"github.com/Nikola-Milovic/vyking-interview/internal/testutil"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:lts",
		mysql.WithDatabase("player_activity_test"),
		mysql.WithUsername("test_user"),
		mysql.WithPassword("test_pass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("port: 3306  MySQL Community Server").
				WithOccurrence(1).
				WithStartupTimeout(30*time.Second),
		),
	)
	require.NoError(t, err)

	host, err := mysqlContainer.Host(ctx)
	require.NoError(t, err)

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	require.NoError(t, err)

	dsn := "test_user:test_pass@tcp(" + host + ":" + port.Port() + ")/player_activity_test?parseTime=true&multiStatements=true"
	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)

	// Run migrations
	migrationsPath := filepath.Join("..", "..", "migrations")
	err = testutil.RunMigrations(db, migrationsPath)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(mysqlContainer); err != nil {
			t.Logf("failed to terminate container: %s", err)
		}
	}

	return db, cleanup
}

func TestService_GetCountryPlayerStats(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := store.New(db)
	mockCountryClient := mock.NewMockCountryAPIClient(ctrl)
	svc := service.New(store, mockCountryClient)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), "RS").
		Return(domain.CountryInfo{
			Name:    "Serbia",
			Region:  "Europe",
			Borders: []string{"BA", "BG", "HR", "HU", "XK", "MK", "ME", "RO"},
		}, nil).
		Times(1)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), "DE").
		Return(domain.CountryInfo{
			Name:    "Germany",
			Region:  "Europe",
			Borders: []string{"AT", "BE", "CZ", "DK", "FR", "LU", "NL", "PL", "CH"},
		}, nil).
		Times(1)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), "BR").
		Return(domain.CountryInfo{
			Name:    "Brazil",
			Region:  "Americas",
			Borders: []string{"AR", "BO", "CO", "GF", "GY", "PY", "PE", "SR", "UY", "VE"},
		}, nil).
		Times(1)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), "UK").
		Return(domain.CountryInfo{
			Name:    "United Kingdom",
			Region:  "Europe",
			Borders: []string{"IE"},
		}, nil).
		Times(1)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), "ES").
		Return(domain.CountryInfo{
			Name:    "Spain",
			Region:  "Europe",
			Borders: []string{"AD", "FR", "GI", "PT", "MA"},
		}, nil).
		Times(1)

	ctx := context.Background()
	req := domain.GetCountryPlayerStatsRequest{
		Limit: 5,
	}

	resp, err := svc.GetCountryPlayerStats(ctx, req)
	require.NoError(t, err)

	assert.Len(t, resp.Stats, 5)

	firstStat := resp.Stats[0]
	assert.Equal(t, "RS", firstStat.CountryCode)
	assert.Equal(t, 10, firstStat.PlayerCount)
	assert.Greater(t, firstStat.TotalBets, 0.0)
	assert.Greater(t, firstStat.AvgBetPerPlayer, 0.0)
	assert.NotNil(t, firstStat.CountryInfo)
	assert.Equal(t, "Serbia", firstStat.CountryInfo.Name)
	assert.Equal(t, "Europe", firstStat.CountryInfo.Region)
	assert.Len(t, firstStat.CountryInfo.Borders, 8)
}

func TestService_GetCountryPlayerStats_WithAPIError(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := store.New(db)
	mockCountryClient := mock.NewMockCountryAPIClient(ctrl)
	svc := service.New(store, mockCountryClient)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), "RS").
		Return(domain.CountryInfo{
			Name:    "Serbia",
			Region:  "Europe",
			Borders: []string{"BA", "BG", "HR", "HU", "XK", "MK", "ME", "RO"},
		}, nil).
		Times(1)

	mockCountryClient.EXPECT().
		GetCountryInfo(gomock.Any(), gomock.Not("RS")).
		Return(domain.CountryInfo{}, fmt.Errorf("API error")).
		AnyTimes()

	ctx := context.Background()
	req := domain.GetCountryPlayerStatsRequest{
		Limit: 3,
	}

	resp, err := svc.GetCountryPlayerStats(ctx, req)
	require.NoError(t, err)

	assert.Len(t, resp.Stats, 3)

	assert.False(t, resp.Stats[0].CountryInfo.IsZero())
	assert.True(t, resp.Stats[1].CountryInfo.IsZero())
	assert.True(t, resp.Stats[2].CountryInfo.IsZero())
}
