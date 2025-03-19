package customTypes

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// uuidFromId custom string type.

var _ basetypes.StringTypable = UuidFromIdStringType{}

type UuidFromIdStringType struct {
	basetypes.StringType
}

func (t UuidFromIdStringType) Equal(o attr.Type) bool {
	other, ok := o.(UuidFromIdStringType)

	if !ok {
		return false
	}

	return t.StringType.Equal(other.StringType)
}

func (t UuidFromIdStringType) String() string {
	return "UuidFromIdStringType"
}

func (t UuidFromIdStringType) ValueFromString(ctx context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	value := UuidFromIdStringValue{
		StringValue: in,
	}

	return value, nil
}

func (t UuidFromIdStringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	attrValue, err := t.StringType.ValueFromTerraform(ctx, in)

	if err != nil {
		return nil, err
	}

	stringValue, ok := attrValue.(basetypes.StringValue)

	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrValue)
	}

	stringValuable, diags := t.ValueFromString(ctx, stringValue)

	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting StringValue to StringValuable: %v", diags)
	}

	return stringValuable, nil
}

func (t UuidFromIdStringType) ValueType(ctx context.Context) attr.Value {
	return UuidFromIdStringValue{}
}

// UuidFromId custom string value.

var _ basetypes.StringValuableWithSemanticEquals = UuidFromIdStringValue{}

type UuidFromIdStringValue struct {
	basetypes.StringValue
}

func (v UuidFromIdStringValue) Equal(o attr.Value) bool {
	other, ok := o.(UuidFromIdStringValue)

	if !ok {
		return false
	}

	return v.StringValue.Equal(other.StringValue)
}

func (v UuidFromIdStringValue) Type(ctx context.Context) attr.Type {
	return UuidFromIdStringType{}
}

func (v UuidFromIdStringValue) StringSemanticEquals(ctx context.Context, newValuable basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(UuidFromIdStringValue)

	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}
	tflog.Debug(ctx, fmt.Sprintf("DEBUG LH Original: %v | NewValue: %v", v.StringValue, newValue.StringValue))
	tflog.Debug(ctx, fmt.Sprintf("DEBUG LH Original Trimmed: %v | NewValue Trimmed: %v", GetUuidStringValueFromId(v.StringValue), GetUuidStringValueFromId(newValue.StringValue)))

	priorStringValue := GetUuidStringValueFromId(v.StringValue)
	newStringValue := GetUuidStringValueFromId(newValue.StringValue)

	return priorStringValue.Equal(newStringValue), diags
}

// String returns a human-readable representation of the String value. Use
// the ValueString method for Terraform data handling instead.
//
// The string returned here is not protected by any compatibility guarantees,
// and is intended for logging and error reporting.
func (v UuidFromIdStringValue) String() string {
	if v.IsUnknown() {
		return attr.UnknownValueString
	}

	if v.IsNull() {
		return attr.NullValueString
	}

	str := v.StringValue.String()
	if strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"") {
		str = strings.TrimSuffix(strings.TrimPrefix(str, "\""), "\"")
	}

	return fmt.Sprintf("%q", GetUuidStringFromId(str))
}

func GetUuidStringValueFromId(value basetypes.StringValue) basetypes.StringValue {
	return basetypes.NewStringValue(GetUuidStringFromId(value.ValueString()))
}

func GetUuidStringFromId(value string) string {
	stringValue := value
	if strings.Contains(value, "/") {
		splitId := strings.Split(value, "/")
		stringValue = splitId[len(splitId)-1]
	}

	return stringValue
}

func NewUuidFromIdStringNull() UuidFromIdStringValue {
	return UuidFromIdStringValue{
		StringValue: basetypes.NewStringNull(),
	}
}

func NewUuidFromIdStringUnknown() UuidFromIdStringValue {
	return UuidFromIdStringValue{
		StringValue: basetypes.NewStringUnknown(),
	}
}

func NewUuidFromIdStringValue(value string) UuidFromIdStringValue {
	return UuidFromIdStringValue{
		StringValue: basetypes.NewStringValue(value),
	}
}

func NewUuidFromIdStringPointerValue(value *string) UuidFromIdStringValue {
	return UuidFromIdStringValue{
		StringValue: basetypes.NewStringPointerValue(value),
	}
}
