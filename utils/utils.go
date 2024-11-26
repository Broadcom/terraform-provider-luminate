package utils

import (
	sdk "bitbucket.org/accezz-io/api-documentation/go/sdk"
	"crypto/md5"
	"fmt"
	"github.com/Broadcom/terraform-provider-luminate/service/utils"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"io"
	"log"
	"regexp"
)

const (
	MAX_SITE_NAME_LENGTH   = 700
	MAX_APP_NAME_LENGTH    = 40
	MAX_POLICY_NAME_LENGTH = 255
	LocalIdpId             = "local"
	DefaultCollection      = "7cef2ccc-ed3e-4812-9ef2-b986c5dac2a5"
	RootCollection         = "6b21619f-f505-41ec-af1b-09350be40000"
	DefaultRDPPort         = "3389"
)

func StringMD5(in string) string {
	h := md5.New()
	io.WriteString(h, in)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func StringInSlice(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ValidateString(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	_, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected to be string"))
		return warns, errs
	}
	return warns, errs
}

func ValidateBool(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	_, ok := v.(bool)
	if !ok {
		errs = append(errs, fmt.Errorf("expected to be bool"))
		return warns, errs
	}
	return warns, errs
}

func ValidateEmail(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string

	userEmail, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("expected to be string"))
		return warns, errs
	}

	if !govalidator.IsEmail(userEmail) {
		errs = append(errs, fmt.Errorf("specified username '%s' is not a valid email address", userEmail))
	}

	return warns, errs
}

func ValidateUuid(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string

	uuid, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("UUID is expected to be string"))
		return warns, errs
	}

	if !govalidator.IsUUID(uuid) {
		errs = append(errs, fmt.Errorf("invalid %s: '%s'", k, uuid))
	}

	return warns, errs
}

func ValidateSiteName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string

	err := ValidateStringPropertyLength(v, MAX_SITE_NAME_LENGTH)
	if err != nil {
		errs = append(errs, err)
	}
	return warns, errs
}

func ValidatePolicyName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string

	err := ValidateStringPropertyLength(v, MAX_POLICY_NAME_LENGTH)
	if err != nil {
		errs = append(errs, err)
	}

	return warns, errs
}

func ValidateApplicationName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string

	err := ValidateStringPropertyLength(v, MAX_APP_NAME_LENGTH)
	if err != nil {
		errs = append(errs, err)
	}

	return warns, errs
}

func ValidateStringPropertyLength(v interface{}, maxLength int) error {
	stringValue, ok := v.(string)
	if !ok {
		return fmt.Errorf("expected to be string")
	}
	if len(stringValue) > maxLength {
		return fmt.Errorf("expected to be string of length %d", MAX_APP_NAME_LENGTH)
	}
	return nil
}

func ParseStringList(stringListInterface []interface{}) []string {
	var stringList []string
	for _, str := range stringListInterface {
		stringList = append(stringList, str.(string))
	}

	return stringList
}

func ValidateTenantRole(role string) bool {
	roles := []string{utils.TenantAdmin, utils.TenantViewer}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func ValidateCollectionRole(role string) bool {
	roles := []string{utils.PolicyOwner, utils.ApplicationOwner}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func ValidateSiteRole(role string) bool {
	roles := []string{utils.SiteEditor, utils.SiteConnectorDeployer}
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false

}

func ValidateEntityType(entityType string) bool {
	types := []string{"User", "Group", "ApiClient"}
	for _, t := range types {
		if t == entityType {
			return true
		}
	}
	return false
}

func ExtractIPAndPort(ipString string) (string, string) {
	pattern := `(?i)(?P<protocol>[\w]+)://(?P<ip>[^\s:]+)(?::(?P<port>\d+))?|(?P<ipOnly>[^\s:]+)`
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(ipString)
	if matches == nil {
		return "", ""
	}
	subMatchesMap := make(map[string]string)
	for i, name := range regex.SubexpNames() {
		if name != "" {
			subMatchesMap[name] = matches[i]
		}
	}
	port := subMatchesMap["port"]
	ip := subMatchesMap["ip"]
	ipOnly := subMatchesMap["ipOnly"]

	if ip == "" {
		return ipOnly, port
	}

	return ip, port
}

func ParseSwaggerError(err error) error {
	e := err.(sdk.GenericSwaggerError)
	model := e.Model().(sdk.ModelApiResponse)
	log.Println(model.Status)
	return errors.Wrap(err, model.Message)
}
