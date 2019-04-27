package integration_test

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wesdean/story-book-api/controllers"
	"github.com/wesdean/story-book-api/database/models"
	"github.com/wesdean/story-book-api/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestForks(t *testing.T) {
	t.Run("GET /forks", func(t *testing.T) {
		seedDb()

		var baseUrl = config.IntegrationTest.ApiUrl + "/forks"

		t.Run("Without body", func(t *testing.T) {
			forksController_Index := func(t *testing.T, token string) {
				t.Run("Default returns all top-level forks", func(t *testing.T) {
					req, err := http.NewRequest("GET", baseUrl, nil)
					if err != nil {
						t.Error(err)
						return
					}
					req.Header.Set("Authorization", token)
					resp, err := netClient.Do(req)
					if err != nil {
						t.Error(err)
						return
					}

					body, _ := ioutil.ReadAll(resp.Body)
					bodyStr := strings.Trim(string(body), "\n")

					if resp.StatusCode != http.StatusOK {
						t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
						return
					}

					var forksResp controllers.ForksControllerForksResponse
					err = json.Unmarshal(body, &forksResp)
					if err != nil {
						t.Error(err)
						return
					}

					expected := 6
					if len(forksResp.Forks) != expected {
						t.Errorf("expected %v, got %v", expected, len(forksResp.Forks))
						return
					}
				})

				t.Run("Get forks by id", func(t *testing.T) {
					req, err := http.NewRequest("GET", baseUrl+"?id=3", nil)
					if err != nil {
						t.Error(err)
						return
					}
					req.Header.Set("Authorization", token)
					resp, err := netClient.Do(req)
					if err != nil {
						t.Error(err)
						return
					}

					body, _ := ioutil.ReadAll(resp.Body)
					bodyStr := strings.Trim(string(body), "\n")

					if resp.StatusCode != http.StatusOK {
						t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
						return
					}

					var forksResp controllers.ForksControllerForksResponse
					err = json.Unmarshal(body, &forksResp)
					if err != nil {
						t.Error(err)
						return
					}

					expectedCount := 1
					if len(forksResp.Forks) != expectedCount {
						t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
						return
					}

					expected := "Test Fork 2"
					if forksResp.Forks[0].Title != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
						return
					}
				})

				t.Run("Get forks by parent", func(t *testing.T) {
					req, err := http.NewRequest("GET", baseUrl+"?parent_id=1", nil)
					if err != nil {
						t.Error(err)
						return
					}
					req.Header.Set("Authorization", token)
					resp, err := netClient.Do(req)
					if err != nil {
						t.Error(err)
						return
					}

					body, _ := ioutil.ReadAll(resp.Body)
					bodyStr := strings.Trim(string(body), "\n")

					if resp.StatusCode != http.StatusOK {
						t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
						return
					}

					var forksResp controllers.ForksControllerForksResponse
					err = json.Unmarshal(body, &forksResp)
					if err != nil {
						t.Error(err)
						return
					}

					expectedCount := 2
					if len(forksResp.Forks) != expectedCount {
						t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
						return
					}

					expected := "Test Fork 1"
					if forksResp.Forks[0].Title != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
						return
					}
				})

				t.Run("Get forks by creator", func(t *testing.T) {
					req, err := http.NewRequest("GET", baseUrl+"?creator_id=2", nil)
					if err != nil {
						t.Error(err)
						return
					}
					req.Header.Set("Authorization", token)
					resp, err := netClient.Do(req)
					if err != nil {
						t.Error(err)
						return
					}

					body, _ := ioutil.ReadAll(resp.Body)
					bodyStr := strings.Trim(string(body), "\n")

					if resp.StatusCode != http.StatusOK {
						t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
						return
					}

					var forksResp controllers.ForksControllerForksResponse
					err = json.Unmarshal(body, &forksResp)
					if err != nil {
						t.Error(err)
						return
					}

					expectedCount := 4
					if len(forksResp.Forks) != expectedCount {
						t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
						return
					}

					expected := "Test Story 2"
					if forksResp.Forks[0].Title != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Title)
						return
					}
				})

				t.Run("Get forks by title", func(t *testing.T) {
					req, err := http.NewRequest("GET", baseUrl+"?title=story", nil)
					if err != nil {
						t.Error(err)
						return
					}
					req.Header.Set("Authorization", token)
					resp, err := netClient.Do(req)
					if err != nil {
						t.Error(err)
						return
					}

					body, _ := ioutil.ReadAll(resp.Body)
					bodyStr := strings.Trim(string(body), "\n")

					if resp.StatusCode != http.StatusOK {
						t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
						return
					}

					var forksResp controllers.ForksControllerForksResponse
					err = json.Unmarshal(body, &forksResp)
					if err != nil {
						t.Error(err)
						return
					}

					expectedCount := 6
					if len(forksResp.Forks) != expectedCount {
						t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
						return
					}

					expected := 1
					if forksResp.Forks[0].Id != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
						return
					}
				})

				t.Run("Get forks by description", func(t *testing.T) {
					req, err := http.NewRequest("GET", baseUrl+"?description=girl", nil)
					if err != nil {
						t.Error(err)
						return
					}
					req.Header.Set("Authorization", token)
					resp, err := netClient.Do(req)
					if err != nil {
						t.Error(err)
						return
					}

					body, _ := ioutil.ReadAll(resp.Body)
					bodyStr := strings.Trim(string(body), "\n")

					if resp.StatusCode != http.StatusOK {
						t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
						return
					}

					var forksResp controllers.ForksControllerForksResponse
					err = json.Unmarshal(body, &forksResp)
					if err != nil {
						t.Error(err)
						return
					}

					expectedCount := 1
					if len(forksResp.Forks) != expectedCount {
						t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
						return
					}

					expected := 1
					if forksResp.Forks[0].Id != expected {
						t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
						return
					}
				})

				t.Run("Get forks by whether they are published", func(t *testing.T) {
					t.Run("Is published", func(t *testing.T) {
						req, err := http.NewRequest("GET", baseUrl+"?is_published=true", nil)
						if err != nil {
							t.Error(err)
							return
						}
						req.Header.Set("Authorization", token)
						resp, err := netClient.Do(req)
						if err != nil {
							t.Error(err)
							return
						}

						body, _ := ioutil.ReadAll(resp.Body)
						bodyStr := strings.Trim(string(body), "\n")

						if resp.StatusCode != http.StatusOK {
							t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
							return
						}

						var forksResp controllers.ForksControllerForksResponse
						err = json.Unmarshal(body, &forksResp)
						if err != nil {
							t.Error(err)
							return
						}

						expectedCount := 2
						if len(forksResp.Forks) != expectedCount {
							t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
							return
						}

						expected := 4
						if forksResp.Forks[0].Id != expected {
							t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
							return
						}
					})

					t.Run("Is not published", func(t *testing.T) {
						req, err := http.NewRequest("GET", baseUrl+"?is_published=false", nil)
						if err != nil {
							t.Error(err)
							return
						}
						req.Header.Set("Authorization", token)
						resp, err := netClient.Do(req)
						if err != nil {
							t.Error(err)
							return
						}

						body, _ := ioutil.ReadAll(resp.Body)
						bodyStr := strings.Trim(string(body), "\n")

						if resp.StatusCode != http.StatusOK {
							t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
							return
						}

						var forksResp controllers.ForksControllerForksResponse
						err = json.Unmarshal(body, &forksResp)
						if err != nil {
							t.Error(err)
							return
						}

						expectedCount := 4
						if len(forksResp.Forks) != expectedCount {
							t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
							return
						}

						expected := 1
						if forksResp.Forks[0].Id != expected {
							t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
							return
						}
					})
				})

				t.Run("Get forks by when they were published", func(t *testing.T) {
					t.Run("Start date only", func(t *testing.T) {
						req, err := http.NewRequest("GET",
							fmt.Sprintf("%s?published_start=%s", baseUrl, url.QueryEscape("2019-03-01 00:00:00-0600")),
							nil)
						if err != nil {
							t.Error(err)
							return
						}
						req.Header.Set("Authorization", token)
						resp, err := netClient.Do(req)
						if err != nil {
							t.Error(err)
							return
						}

						body, _ := ioutil.ReadAll(resp.Body)
						bodyStr := strings.Trim(string(body), "\n")

						if resp.StatusCode != http.StatusOK {
							t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
							return
						}

						var forksResp controllers.ForksControllerForksResponse
						err = json.Unmarshal(body, &forksResp)
						if err != nil {
							t.Error(err)
							return
						}

						expectedCount := 1
						if len(forksResp.Forks) != expectedCount {
							t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
							return
						}

						expected := 4
						if forksResp.Forks[0].Id != expected {
							t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
							return
						}
					})

					t.Run("End date only", func(t *testing.T) {
						req, err := http.NewRequest("GET",
							fmt.Sprintf("%s?published_end=%s", baseUrl, url.QueryEscape("2019-04-27 00:00:00-0600")),
							nil)
						if err != nil {
							t.Error(err)
							return
						}
						req.Header.Set("Authorization", token)
						resp, err := netClient.Do(req)
						if err != nil {
							t.Error(err)
							return
						}

						body, _ := ioutil.ReadAll(resp.Body)
						bodyStr := strings.Trim(string(body), "\n")

						if resp.StatusCode != http.StatusOK {
							t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
							return
						}

						var forksResp controllers.ForksControllerForksResponse
						err = json.Unmarshal(body, &forksResp)
						if err != nil {
							t.Error(err)
							return
						}

						expectedCount := 2
						if len(forksResp.Forks) != expectedCount {
							t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
							return
						}

						expected := 4
						if forksResp.Forks[0].Id != expected {
							t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
							return
						}
					})

					t.Run("Start and end date only", func(t *testing.T) {
						req, err := http.NewRequest("GET",
							fmt.Sprintf("%s?published_start=%s&published_end=%s",
								baseUrl,
								url.QueryEscape("2019-03-01 00:00:00-0600"),
								url.QueryEscape("2019-04-27 00:00:00-0600")),
							nil)
						if err != nil {
							t.Error(err)
							return
						}
						req.Header.Set("Authorization", token)
						resp, err := netClient.Do(req)
						if err != nil {
							t.Error(err)
							return
						}

						body, _ := ioutil.ReadAll(resp.Body)
						bodyStr := strings.Trim(string(body), "\n")

						if resp.StatusCode != http.StatusOK {
							t.Errorf("expected %v, got %v\n%v", http.StatusOK, resp.StatusCode, bodyStr)
							return
						}

						var forksResp controllers.ForksControllerForksResponse
						err = json.Unmarshal(body, &forksResp)
						if err != nil {
							t.Error(err)
							return
						}

						expectedCount := 1
						if len(forksResp.Forks) != expectedCount {
							t.Errorf("expected %v, got %v", expectedCount, len(forksResp.Forks))
							return
						}

						expected := 4
						if forksResp.Forks[0].Id != expected {
							t.Errorf("expected %v, got %v", expected, forksResp.Forks[0].Id)
							return
						}
					})
				})
			}

			t.Run("Unauthenticated user denied access", func(t *testing.T) {
				req, err := http.NewRequest("GET", baseUrl, nil)
				if err != nil {
					t.Error(err)
					return
				}
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})

			t.Run("As superuser", func(t *testing.T) {
				token, err := utils.CreateJWTToken(
					jwt.MapClaims{"user_id": 1},
					[]byte(os.Getenv("AUTH_SECRET")),
				)
				if err != nil {
					t.Error(err)
					return
				}

				forksController_Index(t, token)
			})
		})
	})

	t.Run("POST /forks", func(t *testing.T) {
		var baseUrl = config.IntegrationTest.ApiUrl + "/forks"

		t.Run("Unauthenticated user denied access", func(t *testing.T) {
			seedDb()

			req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"title": "New Fork", "description": "This is a new fork"}`))
			if err != nil {
				t.Error(err)
				return
			}
			resp, err := netClient.Do(req)
			if err != nil {
				t.Error(err)
				return
			}

			body, _ := ioutil.ReadAll(resp.Body)
			bodyStr := strings.Trim(string(body), "\n")

			if resp.StatusCode != http.StatusUnauthorized {
				t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
				return
			}
		})

		t.Run("As superuser", func(t *testing.T) {
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": 1},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Run("Successful fork creation", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"title": "New Fork", "description": "This is a new fork"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				db := openDB()
				defer closeDB(db)
				forkStore := models.NewForkStore(db, nil)
				fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(10))
				if err != nil {
					t.Error(err)
					return
				}

				if fork == nil {
					t.Error("expected fork, got nil")
					return
				}
			})

			t.Run("Empty request body returns error", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, nil)
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusBadRequest {
					t.Errorf("expected %v, got %v\n%v", http.StatusBadRequest, resp.StatusCode, bodyStr)
					return
				}

				expected := `{"error":"invalid request body"}`
				if bodyStr != expected {
					t.Errorf("expected %v, got %v", expected, bodyStr)
					return
				}
			})

		})

		t.Run("As owner", func(t *testing.T) {
			userId := 2
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": userId},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Run("Can create top-level fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				db := openDB()
				defer closeDB(db)
				forkStore := models.NewForkStore(db, nil)
				fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(10))
				if err != nil {
					t.Error(err)
					return
				}

				if fork == nil {
					t.Error("expected fork, got nil")
					return
				}
			})

			t.Run("Can create fork from owned fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				var forkResp controllers.ForksControllerForkResponse
				err = json.Unmarshal(body, &forkResp)
				if err != nil {
					t.Error(err)
					return
				}

				expectedId := 10
				if forkResp.Fork.Id != expectedId {
					t.Errorf("expected %v, got %v", expectedId, forkResp.Fork.Id)
					return
				}
			})

			t.Run("Cannot create fork from unowned fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 8, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})
		})

		t.Run("As author", func(t *testing.T) {
			userId := 3
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": userId},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Run("Can create top-level fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				db := openDB()
				defer closeDB(db)
				forkStore := models.NewForkStore(db, nil)
				fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(10))
				if err != nil {
					t.Error(err)
					return
				}

				if fork == nil {
					t.Error("expected fork, got nil")
					return
				}
			})

			t.Run("Can create fork from owned fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 2, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				var forkResp controllers.ForksControllerForkResponse
				err = json.Unmarshal(body, &forkResp)
				if err != nil {
					t.Error(err)
					return
				}

				expectedId := 10
				if forkResp.Fork.Id != expectedId {
					t.Errorf("expected %v, got %v", expectedId, forkResp.Fork.Id)
					return
				}
			})

			t.Run("Cannot create fork from unowned fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})
		})

		t.Run("As editor", func(t *testing.T) {
			userId := 4
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": userId},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Run("Can create top-level fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				db := openDB()
				defer closeDB(db)
				forkStore := models.NewForkStore(db, nil)
				fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(10))
				if err != nil {
					t.Error(err)
					return
				}

				if fork == nil {
					t.Error("expected fork, got nil")
					return
				}
			})

			t.Run("Cannot create fork from edited fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})

			t.Run("Cannot create fork from unowned fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})
		})

		t.Run("As proofreader", func(t *testing.T) {
			userId := 5
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": userId},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Run("Can create top-level fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				db := openDB()
				defer closeDB(db)
				forkStore := models.NewForkStore(db, nil)
				fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(10))
				if err != nil {
					t.Error(err)
					return
				}

				if fork == nil {
					t.Error("expected fork, got nil")
					return
				}
			})

			t.Run("Cannot create fork from proofread fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})

			t.Run("Cannot create fork from unproofread fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})
		})

		t.Run("As reader", func(t *testing.T) {
			userId := 6
			token, err := utils.CreateJWTToken(
				jwt.MapClaims{"user_id": userId},
				[]byte(os.Getenv("AUTH_SECRET")),
			)
			if err != nil {
				t.Error(err)
				return
			}

			t.Run("Can create top-level fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusCreated {
					t.Errorf("expected %v, got %v\n%v", http.StatusCreated, resp.StatusCode, bodyStr)
					return
				}

				db := openDB()
				defer closeDB(db)
				forkStore := models.NewForkStore(db, nil)
				fork, err := forkStore.GetFork(models.NewForkQueryOptions().Id(10))
				if err != nil {
					t.Error(err)
					return
				}

				if fork == nil {
					t.Error("expected fork, got nil")
					return
				}
			})

			t.Run("Cannot create fork from read fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 4, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})

			t.Run("Cannot create fork from unread fork", func(t *testing.T) {
				seedDb()

				req, err := http.NewRequest("POST", baseUrl, strings.NewReader(`{"ParentId": 7, "Title": "Newly Created Fork", "Description": "This is a fork!"}`))
				if err != nil {
					t.Error(err)
					return
				}
				req.Header.Set("Authorization", token)
				resp, err := netClient.Do(req)
				if err != nil {
					t.Error(err)
					return
				}

				body, _ := ioutil.ReadAll(resp.Body)
				bodyStr := strings.Trim(string(body), "\n")

				if resp.StatusCode != http.StatusUnauthorized {
					t.Errorf("expected %v, got %v\n%v", http.StatusUnauthorized, resp.StatusCode, bodyStr)
					return
				}
			})
		})
	})
	//todo Update fork tests
	//todo Delete fork tests
}
