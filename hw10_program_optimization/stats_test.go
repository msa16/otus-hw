//go:build !bench
// +build !bench

package hw10programoptimization

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDomainStat(t *testing.T) {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	t.Run("find 'com'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{"browsedrive.gov": 1}, result)
	})

	t.Run("find 'unknown'", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "unknown")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})
}

func TestGetDomainStatEmptyData(t *testing.T) {
	data := `{"Id":1, "Email":""}
{"Id":2,"Email":"mLynch@broWsecat.com"}
{"Id":3,"Email":"RoseSmith@Browsecat.com"}
{"Id":4}
{"Id":5,"Email":"nulla@Linktype.com"}`

	t.Run("find 'com' with absent email", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"browsecat.com": 2,
			"linktype.com":  1,
		}, result)
	})

	t.Run("find 'gov' in empty data", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(""), "gov")
		require.NoError(t, err)
		require.Equal(t, DomainStat{}, result)
	})

}

func TestGetDomainStatManyDots(t *testing.T) {
	data := `{"Id":1,"Email":"ivanoff@int.company.com"}
{"Id":2,"Email":"alex.petroff@int.company.com"}
{"Id":3,"Email":"J.R.RoseSmith@company.com"}
{"Id":4,"Email":"nulla@ext.company.com"}
{"Id":5,"Email":"a.sveta.kurnikova@ext.company.com"}`
	t.Run("many dots to left and right from @", func(t *testing.T) {
		result, err := GetDomainStat(bytes.NewBufferString(data), "com")
		require.NoError(t, err)
		require.Equal(t, DomainStat{
			"int.company.com": 2,
			"company.com":     1,
			"ext.company.com": 2,
		}, result)
	})
}
