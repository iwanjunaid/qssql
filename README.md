# QSSQL

## Example

```go
package main

import (
	"fmt"
	"log"
	_ "reflect"

	iface "github.com/iwanjunaid/qssql/pkg/interfaces"
)

func main() {
	template := `
	SELECT *
	FROM books
	{{.Where}}
	{{.OrderBy}}
	GROUP BY id
	{{.Offset}}
	{{.Limit}}
	`

	x := "title=Hello" +
		"&price=50000" +
		"&customer_price[gt]=51000" +
		"&customer_price[lt]=65000" +
		"&student_price[gte]=51000" +
		"&student_price[lte]=60000" +
		"&author[]=deitel&author[]=linus" +
		"&editor[like]=John" +
		"&editor[hlike]=Jenna" +
		"&a[hilike]=HILIKERobert" +
		"&b[hne]=HNERobert" +
		"&c[hin]=HINRobert" +
		"&d[hnin]=HNINRobert" +
		"&e[hgt]=100" +
		"&f[hgte]=105" +
		"&g[hlt]=200" +
		"&h[hlte]=250" +
		"&publisher[in][]=oreally&publisher[in][]=oreilly" +
		"&isbn[nin][]=123" +
		"&limit=10" +
		"&offset=1&sort[]=title&sort[]=created_at"

	qssql := New(x, template, func(q iface.GenericQSSQL) {
		q.AddWhereValues([]string{"t1.isbn = '123'"})
		q.AddWhereAliases(map[string]string{
			"title": "t1.title",
			"price": "t1.price",
		})
	})

	err := qssql.Parse()

	if err != nil {
		log.Fatal(err)
	}

	sql, err := qssql.GetSQL()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("SQL =>", sql)
	fmt.Println("Bind values =>", qssql.GetBindValues())
}
```