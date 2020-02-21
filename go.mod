module bitbucket.org/accezz-io/terraform-provider-symcsc

// At the moment of writing git.apache.org is down. Using github repo to make build work
replace git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999

require (
	bitbucket.org/accezz-io/api-documentation v0.0.0-20200211094502-8f6a88172b30db8a7bc6aebece21932b282337f2
	github.com/antihax/optional v0.0.0-20180407024304-ca021399b1a6
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a
	github.com/hashicorp/terraform v0.12.2
	github.com/pkg/errors v0.0.0-20170505043639-c605e284fe17
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.3.0
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
)

go 1.13
