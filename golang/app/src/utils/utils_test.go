package utils

import (
	"net/url"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	//run the tests
	exitVal := m.Run()

	os.Exit(exitVal)
}

type dateFromQueryTest struct {
	queryKey     string
	query        url.Values
	expectedDate time.Time
	expectedErr  bool
}

func TestDateFromQuery(t *testing.T) {
	query := url.Values{
		"date1":        []string{"2002-01-25T00:00:00+00:00"},
		"date2":        []string{"2022-01-17T07:32:30+00:00"},
		"date3":        []string{"2022-01-17T07:32:30Z"},
		"invalid_date": []string{"foo"},
	}

	dateOne := time.Date(2002, 01, 25, 0, 0, 0, 0, time.UTC)
	dateFull := time.Date(2022, 01, 17, 7, 32, 30, 0, time.UTC)

	dateFromQueryTests := []dateFromQueryTest{
		//test normal behaviors
		{
			queryKey:     "date1",
			query:        query,
			expectedDate: dateOne,
			expectedErr:  false,
		},
		{
			queryKey:     "date2",
			query:        query,
			expectedDate: dateFull,
			expectedErr:  false,
		},
		{
			queryKey:     "date3",
			query:        query,
			expectedDate: dateFull,
			expectedErr:  false,
		},
		//test missing query key
		{
			queryKey:     "missing_date",
			query:        query,
			expectedDate: time.Time{},
			expectedErr:  false,
		},
		//test invalid value
		{
			queryKey:     "invalid_date",
			query:        query,
			expectedDate: time.Time{},
			expectedErr:  true,
		},
	}

	for _, item := range dateFromQueryTests {
		resultDate, resultErr := DateFromQuery(item.queryKey, item.query)
		if !item.expectedErr && resultErr != nil {
			t.Errorf("DateFromQuery for %v output err %v not expected", item.queryKey, resultErr.Error())
		}
		if item.expectedErr && resultErr == nil {
			t.Errorf("DateFromQuery for %v output err expected but not found", item.queryKey)
		}
		if !resultDate.Equal(item.expectedDate) {
			t.Errorf("DateFromQuery for %v output %v but expected %v", item.queryKey, resultDate, item.expectedDate)
		}
	}
}

type stringFromQueryTest struct {
	queryKey       string
	query          url.Values
	expectedString string
}

func TestStringFromQuery(t *testing.T) {
	query := url.Values{
		"string1": []string{"string1"},
		"string2": []string{"string2"},
		"string3": []string{"string3"},
	}

	stringFromQueryTests := []stringFromQueryTest{
		//test normal behaviors
		{
			queryKey:       "string1",
			query:          query,
			expectedString: "string1",
		},
		{
			queryKey:       "string2",
			query:          query,
			expectedString: "string2",
		},
		{
			queryKey:       "string3",
			query:          query,
			expectedString: "string3",
		},
		//test missing query key
		{
			queryKey:       "missing_string",
			query:          query,
			expectedString: "",
		},
	}

	for _, item := range stringFromQueryTests {
		resultString := StringFromQuery(item.queryKey, item.query)
		if resultString != item.expectedString {
			t.Errorf("StringFromQuery for %v output %v but expected %v", item.queryKey, resultString, item.expectedString)
		}
	}
}

type int64FromQueryTest struct {
	queryKey      string
	query         url.Values
	expectedInt64 int64
	expectedErr   bool
}

func TestInt64FromQuery(t *testing.T) {
	query := url.Values{
		"int1":        []string{"2"},
		"int2":        []string{"237165"},
		"invalid_int": []string{"foo"},
	}

	int64FromQueryTests := []int64FromQueryTest{
		//test normal behaviors
		{
			queryKey:      "int1",
			query:         query,
			expectedInt64: 2,
			expectedErr:   false,
		},
		{
			queryKey:      "int2",
			query:         query,
			expectedInt64: 237165,
			expectedErr:   false,
		},
		//test missing query key
		{
			queryKey:      "missing_int",
			query:         query,
			expectedInt64: 0,
			expectedErr:   false,
		},
		//test invalid value
		{
			queryKey:      "invalid_int",
			query:         query,
			expectedInt64: 0,
			expectedErr:   true,
		},
	}

	for _, item := range int64FromQueryTests {
		resultInt64, resultErr := Int64FromQuery(item.queryKey, item.query)
		if !item.expectedErr && resultErr != nil {
			t.Errorf("Int64FromQuery for %v output err %v not expected", item.queryKey, resultErr.Error())
		}
		if item.expectedErr && resultErr == nil {
			t.Errorf("Int64FromQuery for %v output err expected but not found", item.queryKey)
		}
		if resultInt64 != item.expectedInt64 {
			t.Errorf("Int64FromQuery for %v output %v but expected %v", item.queryKey, resultInt64, item.expectedInt64)
		}
	}
}
