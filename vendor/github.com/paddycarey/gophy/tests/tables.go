// This file contains all test data/tables for our table driven tests
package tests

var SearchGifsTestData = []struct {
	ApiKey              string
	Q                   string
	Limit               int
	Offset              int
	Rating              string
	ExpectedError       bool
	ExpectedNumReturned int
	ExpectedTotalCount  int
}{
	{"", "funny cats", 10, 0, "", false, 10, 1913},
	{"", "funny cats", 10, -1, "", true, 10, 1913},
	{"", "funny cats", 10, 0, "ad", true, 10, 1913},
	{"", "tom hiddleston", 20, 10, "g", false, 20, 1987},
	{"", "table flip", 99, 0, "", false, 99, 99},
	{"", "facepalm", 100, 0, "", false, 100, 360},
	{"", "picard", 20, 11, "", false, 20, 180},
	{"", "kittens", 30, 32, "", false, 30, 4110},
	{"", "party hard", 4, 0, "", false, 4, 719},
	{"", "friday", 50, 20, "", false, 50, 3237},
	{"", "snow", 60, 30, "", false, 60, 5483},
	{"", "new york", 100, 9000, "", false, 0, 4666},
	{"", "new york", 101, 9000, "", true, 0, 4666},
	{"", "", 100, 9000, "", true, 0, 4666},
	{"abc", "banana", 10, 0, "", true, 10, 1911},
}

var SearchStickersTestData = []struct {
	ApiKey              string
	Q                   string
	Limit               int
	Offset              int
	Rating              string
	ExpectedError       bool
	ExpectedNumReturned int
	ExpectedTotalCount  int
}{
	{"", "funny cats", 10, 0, "", false, 10, 15},
	{"", "funny cats", 10, -1, "", true, 10, 15},
	{"", "funny cats", 10, 0, "ad", true, 10, 15},
	{"", "tom hiddleston", 20, 10, "g", false, 0, 4},
	{"", "table flip", 99, 0, "", false, 0, 0},
	{"", "facepalm", 100, 0, "", false, 21, 21},
	{"", "picard", 20, 11, "", false, 0, 0},
	{"", "kittens", 30, 32, "", false, 30, 100},
	{"", "party hard", 4, 0, "", false, 4, 6},
	{"", "friday", 50, 20, "", false, 30, 50},
	{"", "snow", 60, 30, "", false, 34, 64},
	{"", "new york", 100, 9000, "", false, 0, 25},
	{"", "new york", 101, 9000, "", true, 0, 25},
	{"", "", 100, 9000, "", true, 0, 4666},
	{"abc", "banana", 10, 0, "", true, 10, 1911},
}

var GetGifByIdTestData = []struct {
	ApiKey                  string
	Id                      string
	ExpectedGetGifByIdError bool
}{
	{"", "qe3gu0bH7Shj2", false},
	{"", "uW3OQZo9olSZG", false},
	{"", "3uWOZ9lZ0HSj2", true},
	{"sf", "qe3gu0bH7Shj2", true},
}

var GetGifsByIdTestData = []struct {
	ApiKey                   string
	Ids                      []string
	ExpectedNumReturned      int
	ExpectedGetGifsByIdError bool
}{
	{"", []string{
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "",
		"", "", "", "", "", "", "", "", "", "", "",
	}, 0, true},
	{"", []string{"qe3gu0bH7Shj2", "uW3OQZo9olSZG"}, 2, false},
	{"", []string{"uW3OQZo9olSZG"}, 1, false},
	{"", []string{"3uWOZ9lZ0HSj2"}, 0, false},
	{"sf", []string{"qe3gu0bH7Shj2"}, 0, true},
}

var TranslateGifTestData = []struct {
	ApiKey                 string
	Q                      string
	Rating                 string
	ExpectedTranslateError bool
}{
	{"", "", "", true},
	{"", "batman", "gf", true},
	{"", "superman", "", false},
	{"", "superman", "g", false},
	{"", "burn", "", false},
	{"", "smackdown", "", false},
	{"sf", "facepalm", "", true},
}

var TranslateStickerTestData = []struct {
	ApiKey                 string
	Q                      string
	Rating                 string
	ExpectedTranslateError bool
}{
	{"", "", "", true},
	{"", "batman", "gf", true},
	{"", "superman", "", false},
	{"", "superman", "g", false},
	{"", "burn", "", false},
	{"", "smackdown", "", false},
	{"sf", "facepalm", "", true},
}

var TrendingGifsTestData = []struct {
	ApiKey              string
	Rating              string
	Limit               int
	ExpectedError       bool
	ExpectedNumReturned int
}{
	{"", "", -10, true, 0},
	{"", "", 101, true, 0},
	{"", "", 10, false, 10},
	{"", "g", 100, false, 100},
	{"sdf", "", 100, true, 0},
	{"", "d", 100, true, 0},
}
var TrendingStickersTestData = []struct {
	ApiKey              string
	Rating              string
	Limit               int
	ExpectedError       bool
	ExpectedNumReturned int
}{
	{"", "", -10, true, 0},
	{"", "", 101, true, 0},
	{"", "", 10, false, 10},
	{"", "g", 100, false, 100},
	{"sdf", "", 100, true, 0},
	{"", "d", 100, true, 0},
}
