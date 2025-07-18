// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package expression

import (
	"fmt"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/tidb/pkg/parser/ast"
	"github.com/pingcap/tidb/pkg/parser/mysql"
	"github.com/pingcap/tidb/pkg/types"
	"github.com/pingcap/tidb/pkg/util/codec"
	"github.com/pingcap/tidb/pkg/util/collate"
	"github.com/pingcap/tipb/go-tipb"
)

// PbTypeToFieldType converts tipb.FieldType to FieldType
func PbTypeToFieldType(tp *tipb.FieldType) *types.FieldType {
	ftb := types.NewFieldTypeBuilder()
	ft := ftb.SetType(byte(tp.Tp)).SetFlag(uint(tp.Flag)).SetFlen(int(tp.Flen)).SetDecimal(int(tp.Decimal)).SetCharset(tp.Charset).SetCollate(collate.ProtoToCollation(tp.Collate)).BuildP()
	ft.SetElems(tp.Elems)
	return ft
}

func getSignatureByPB(ctx BuildContext, sigCode tipb.ScalarFuncSig, tp *tipb.FieldType, args []Expression) (f builtinFunc, e error) {
	fieldTp := PbTypeToFieldType(tp)
	base, err := newBaseBuiltinFuncWithFieldType(fieldTp, args)
	if err != nil {
		return nil, err
	}
	maxAllowedPacket := ctx.GetEvalCtx().GetMaxAllowedPacket()
	switch sigCode {
	case tipb.ScalarFuncSig_CastIntAsInt:
		f = &builtinCastIntAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastIntAsReal:
		f = &builtinCastIntAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastIntAsString:
		f = &builtinCastIntAsStringSig{base}
	case tipb.ScalarFuncSig_CastIntAsDecimal:
		f = &builtinCastIntAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastIntAsTime:
		f = &builtinCastIntAsTimeSig{base}
	case tipb.ScalarFuncSig_CastIntAsDuration:
		f = &builtinCastIntAsDurationSig{base}
	case tipb.ScalarFuncSig_CastIntAsJson:
		f = &builtinCastIntAsJSONSig{base}
	case tipb.ScalarFuncSig_CastRealAsInt:
		f = &builtinCastRealAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastRealAsReal:
		f = &builtinCastRealAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastRealAsString:
		f = &builtinCastRealAsStringSig{base}
	case tipb.ScalarFuncSig_CastRealAsDecimal:
		f = &builtinCastRealAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastRealAsTime:
		f = &builtinCastRealAsTimeSig{base}
	case tipb.ScalarFuncSig_CastRealAsDuration:
		f = &builtinCastRealAsDurationSig{base}
	case tipb.ScalarFuncSig_CastRealAsJson:
		f = &builtinCastRealAsJSONSig{base}
	case tipb.ScalarFuncSig_CastDecimalAsInt:
		f = &builtinCastDecimalAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastDecimalAsReal:
		f = &builtinCastDecimalAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastDecimalAsString:
		f = &builtinCastDecimalAsStringSig{base}
	case tipb.ScalarFuncSig_CastDecimalAsDecimal:
		f = &builtinCastDecimalAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastDecimalAsTime:
		f = &builtinCastDecimalAsTimeSig{base}
	case tipb.ScalarFuncSig_CastDecimalAsDuration:
		f = &builtinCastDecimalAsDurationSig{base}
	case tipb.ScalarFuncSig_CastDecimalAsJson:
		f = &builtinCastDecimalAsJSONSig{base}
	case tipb.ScalarFuncSig_CastStringAsInt:
		f = &builtinCastStringAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastStringAsReal:
		f = &builtinCastStringAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastStringAsString:
		f = &builtinCastStringAsStringSig{base}
	case tipb.ScalarFuncSig_CastStringAsDecimal:
		f = &builtinCastStringAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastStringAsTime:
		f = &builtinCastStringAsTimeSig{base}
	case tipb.ScalarFuncSig_CastStringAsDuration:
		f = &builtinCastStringAsDurationSig{base}
	case tipb.ScalarFuncSig_CastStringAsJson:
		f = &builtinCastStringAsJSONSig{base}
	case tipb.ScalarFuncSig_CastTimeAsInt:
		f = &builtinCastTimeAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastTimeAsReal:
		f = &builtinCastTimeAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastTimeAsString:
		f = &builtinCastTimeAsStringSig{base}
	case tipb.ScalarFuncSig_CastTimeAsDecimal:
		f = &builtinCastTimeAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastTimeAsTime:
		f = &builtinCastTimeAsTimeSig{base}
	case tipb.ScalarFuncSig_CastTimeAsDuration:
		f = &builtinCastTimeAsDurationSig{base}
	case tipb.ScalarFuncSig_CastTimeAsJson:
		f = &builtinCastTimeAsJSONSig{base}
	case tipb.ScalarFuncSig_CastDurationAsInt:
		f = &builtinCastDurationAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastDurationAsReal:
		f = &builtinCastDurationAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastDurationAsString:
		f = &builtinCastDurationAsStringSig{base}
	case tipb.ScalarFuncSig_CastDurationAsDecimal:
		f = &builtinCastDurationAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastDurationAsTime:
		f = &builtinCastDurationAsTimeSig{base}
	case tipb.ScalarFuncSig_CastDurationAsDuration:
		f = &builtinCastDurationAsDurationSig{base}
	case tipb.ScalarFuncSig_CastDurationAsJson:
		f = &builtinCastDurationAsJSONSig{base}
	case tipb.ScalarFuncSig_CastJsonAsInt:
		f = &builtinCastJSONAsIntSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastJsonAsReal:
		f = &builtinCastJSONAsRealSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastJsonAsString:
		f = &builtinCastJSONAsStringSig{base}
	case tipb.ScalarFuncSig_CastJsonAsDecimal:
		f = &builtinCastJSONAsDecimalSig{newBaseBuiltinCastFunc(base, false)}
	case tipb.ScalarFuncSig_CastJsonAsTime:
		f = &builtinCastJSONAsTimeSig{base}
	case tipb.ScalarFuncSig_CastJsonAsDuration:
		f = &builtinCastJSONAsDurationSig{base}
	case tipb.ScalarFuncSig_CastJsonAsJson:
		f = &builtinCastJSONAsJSONSig{base}
	case tipb.ScalarFuncSig_CoalesceInt:
		f = &builtinCoalesceIntSig{base}
	case tipb.ScalarFuncSig_CoalesceReal:
		f = &builtinCoalesceRealSig{base}
	case tipb.ScalarFuncSig_CoalesceDecimal:
		f = &builtinCoalesceDecimalSig{base}
	case tipb.ScalarFuncSig_CoalesceString:
		f = &builtinCoalesceStringSig{base}
	case tipb.ScalarFuncSig_CoalesceTime:
		f = &builtinCoalesceTimeSig{base}
	case tipb.ScalarFuncSig_CoalesceDuration:
		f = &builtinCoalesceDurationSig{base}
	case tipb.ScalarFuncSig_CoalesceJson:
		f = &builtinCoalesceJSONSig{base}
	case tipb.ScalarFuncSig_LTInt:
		f = &builtinLTIntSig{base}
	case tipb.ScalarFuncSig_LTReal:
		f = &builtinLTRealSig{base}
	case tipb.ScalarFuncSig_LTDecimal:
		f = &builtinLTDecimalSig{base}
	case tipb.ScalarFuncSig_LTString:
		f = &builtinLTStringSig{base}
	case tipb.ScalarFuncSig_LTTime:
		f = &builtinLTTimeSig{base}
	case tipb.ScalarFuncSig_LTDuration:
		f = &builtinLTDurationSig{base}
	case tipb.ScalarFuncSig_LTJson:
		f = &builtinLTJSONSig{base}
	case tipb.ScalarFuncSig_LEInt:
		f = &builtinLEIntSig{base}
	case tipb.ScalarFuncSig_LEReal:
		f = &builtinLERealSig{base}
	case tipb.ScalarFuncSig_LEDecimal:
		f = &builtinLEDecimalSig{base}
	case tipb.ScalarFuncSig_LEString:
		f = &builtinLEStringSig{base}
	case tipb.ScalarFuncSig_LETime:
		f = &builtinLETimeSig{base}
	case tipb.ScalarFuncSig_LEDuration:
		f = &builtinLEDurationSig{base}
	case tipb.ScalarFuncSig_LEJson:
		f = &builtinLEJSONSig{base}
	case tipb.ScalarFuncSig_GTInt:
		f = &builtinGTIntSig{base}
	case tipb.ScalarFuncSig_GTReal:
		f = &builtinGTRealSig{base}
	case tipb.ScalarFuncSig_GTDecimal:
		f = &builtinGTDecimalSig{base}
	case tipb.ScalarFuncSig_GTString:
		f = &builtinGTStringSig{base}
	case tipb.ScalarFuncSig_GTTime:
		f = &builtinGTTimeSig{base}
	case tipb.ScalarFuncSig_GTDuration:
		f = &builtinGTDurationSig{base}
	case tipb.ScalarFuncSig_GTJson:
		f = &builtinGTJSONSig{base}
	case tipb.ScalarFuncSig_GreatestInt:
		f = &builtinGreatestIntSig{base}
	case tipb.ScalarFuncSig_GreatestReal:
		f = &builtinGreatestRealSig{base}
	case tipb.ScalarFuncSig_GreatestDecimal:
		f = &builtinGreatestDecimalSig{base}
	case tipb.ScalarFuncSig_GreatestString:
		f = &builtinGreatestStringSig{base}
	case tipb.ScalarFuncSig_GreatestTime:
		f = &builtinGreatestTimeSig{base, false}
	case tipb.ScalarFuncSig_GreatestDate:
		f = &builtinGreatestTimeSig{base, true}
	case tipb.ScalarFuncSig_GreatestCmpStringAsTime:
		f = &builtinGreatestCmpStringAsTimeSig{base, false}
	case tipb.ScalarFuncSig_GreatestCmpStringAsDate:
		f = &builtinGreatestCmpStringAsTimeSig{base, true}
	case tipb.ScalarFuncSig_GreatestDuration:
		f = &builtinGreatestDurationSig{base}
	case tipb.ScalarFuncSig_LeastInt:
		f = &builtinLeastIntSig{base}
	case tipb.ScalarFuncSig_LeastReal:
		f = &builtinLeastRealSig{base}
	case tipb.ScalarFuncSig_LeastDecimal:
		f = &builtinLeastDecimalSig{base}
	case tipb.ScalarFuncSig_LeastString:
		f = &builtinLeastStringSig{base}
	case tipb.ScalarFuncSig_LeastTime:
		f = &builtinLeastTimeSig{base, false}
	case tipb.ScalarFuncSig_LeastDate:
		f = &builtinLeastTimeSig{base, true}
	case tipb.ScalarFuncSig_LeastCmpStringAsTime:
		f = &builtinLeastCmpStringAsTimeSig{base, false}
	case tipb.ScalarFuncSig_LeastCmpStringAsDate:
		f = &builtinLeastCmpStringAsTimeSig{base, true}
	case tipb.ScalarFuncSig_LeastDuration:
		f = &builtinLeastDurationSig{base}
	case tipb.ScalarFuncSig_IntervalInt:
		f = &builtinIntervalIntSig{base, false} // Since interval function won't be pushed down to TiKV, therefore it doesn't matter what value we give to hasNullable
	case tipb.ScalarFuncSig_IntervalReal:
		f = &builtinIntervalRealSig{base, false}
	case tipb.ScalarFuncSig_GEInt:
		f = &builtinGEIntSig{base}
	case tipb.ScalarFuncSig_GEReal:
		f = &builtinGERealSig{base}
	case tipb.ScalarFuncSig_GEDecimal:
		f = &builtinGEDecimalSig{base}
	case tipb.ScalarFuncSig_GEString:
		f = &builtinGEStringSig{base}
	case tipb.ScalarFuncSig_GETime:
		f = &builtinGETimeSig{base}
	case tipb.ScalarFuncSig_GEDuration:
		f = &builtinGEDurationSig{base}
	case tipb.ScalarFuncSig_GEJson:
		f = &builtinGEJSONSig{base}
	case tipb.ScalarFuncSig_EQInt:
		f = &builtinEQIntSig{base}
	case tipb.ScalarFuncSig_EQReal:
		f = &builtinEQRealSig{base}
	case tipb.ScalarFuncSig_EQDecimal:
		f = &builtinEQDecimalSig{base}
	case tipb.ScalarFuncSig_EQString:
		f = &builtinEQStringSig{base}
	case tipb.ScalarFuncSig_EQTime:
		f = &builtinEQTimeSig{base}
	case tipb.ScalarFuncSig_EQDuration:
		f = &builtinEQDurationSig{base}
	case tipb.ScalarFuncSig_EQJson:
		f = &builtinEQJSONSig{base}
	case tipb.ScalarFuncSig_NEInt:
		f = &builtinNEIntSig{base}
	case tipb.ScalarFuncSig_NEReal:
		f = &builtinNERealSig{base}
	case tipb.ScalarFuncSig_NEDecimal:
		f = &builtinNEDecimalSig{base}
	case tipb.ScalarFuncSig_NEString:
		f = &builtinNEStringSig{base}
	case tipb.ScalarFuncSig_NETime:
		f = &builtinNETimeSig{base}
	case tipb.ScalarFuncSig_NEDuration:
		f = &builtinNEDurationSig{base}
	case tipb.ScalarFuncSig_NEJson:
		f = &builtinNEJSONSig{base}
	case tipb.ScalarFuncSig_NullEQInt:
		f = &builtinNullEQIntSig{base}
	case tipb.ScalarFuncSig_NullEQReal:
		f = &builtinNullEQRealSig{base}
	case tipb.ScalarFuncSig_NullEQDecimal:
		f = &builtinNullEQDecimalSig{base}
	case tipb.ScalarFuncSig_NullEQString:
		f = &builtinNullEQStringSig{base}
	case tipb.ScalarFuncSig_NullEQTime:
		f = &builtinNullEQTimeSig{base}
	case tipb.ScalarFuncSig_NullEQDuration:
		f = &builtinNullEQDurationSig{base}
	case tipb.ScalarFuncSig_NullEQJson:
		f = &builtinNullEQJSONSig{base}
	case tipb.ScalarFuncSig_PlusReal:
		f = &builtinArithmeticPlusRealSig{base}
	case tipb.ScalarFuncSig_PlusDecimal:
		f = &builtinArithmeticPlusDecimalSig{base}
	case tipb.ScalarFuncSig_PlusInt:
		f = &builtinArithmeticPlusIntSig{base}
	case tipb.ScalarFuncSig_MinusReal:
		f = &builtinArithmeticMinusRealSig{base}
	case tipb.ScalarFuncSig_MinusDecimal:
		f = &builtinArithmeticMinusDecimalSig{base}
	case tipb.ScalarFuncSig_MinusInt:
		f = &builtinArithmeticMinusIntSig{base}
	case tipb.ScalarFuncSig_MultiplyReal:
		f = &builtinArithmeticMultiplyRealSig{base, false}
	case tipb.ScalarFuncSig_MultiplyDecimal:
		f = &builtinArithmeticMultiplyDecimalSig{base}
	case tipb.ScalarFuncSig_MultiplyInt:
		f = &builtinArithmeticMultiplyIntSig{base}
	case tipb.ScalarFuncSig_DivideReal:
		f = &builtinArithmeticDivideRealSig{base}
	case tipb.ScalarFuncSig_DivideDecimal:
		f = &builtinArithmeticDivideDecimalSig{base}
	case tipb.ScalarFuncSig_IntDivideInt:
		f = &builtinArithmeticIntDivideIntSig{base}
	case tipb.ScalarFuncSig_IntDivideDecimal:
		f = &builtinArithmeticIntDivideDecimalSig{base}
	case tipb.ScalarFuncSig_ModReal:
		f = &builtinArithmeticModRealSig{base}
	case tipb.ScalarFuncSig_ModDecimal:
		f = &builtinArithmeticModDecimalSig{base}
	case tipb.ScalarFuncSig_ModIntUnsignedUnsigned:
		f = &builtinArithmeticModIntUnsignedUnsignedSig{base}
	case tipb.ScalarFuncSig_ModIntUnsignedSigned:
		f = &builtinArithmeticModIntUnsignedSignedSig{base}
	case tipb.ScalarFuncSig_ModIntSignedUnsigned:
		f = &builtinArithmeticModIntSignedUnsignedSig{base}
	case tipb.ScalarFuncSig_ModIntSignedSigned:
		f = &builtinArithmeticModIntSignedSignedSig{base}
	case tipb.ScalarFuncSig_MultiplyIntUnsigned:
		f = &builtinArithmeticMultiplyIntUnsignedSig{base}
	case tipb.ScalarFuncSig_AbsInt:
		f = &builtinAbsIntSig{base}
	case tipb.ScalarFuncSig_AbsUInt:
		f = &builtinAbsUIntSig{base}
	case tipb.ScalarFuncSig_AbsReal:
		f = &builtinAbsRealSig{base}
	case tipb.ScalarFuncSig_AbsDecimal:
		f = &builtinAbsDecSig{base}
	case tipb.ScalarFuncSig_CeilIntToDec:
		f = &builtinCeilIntToDecSig{base}
	case tipb.ScalarFuncSig_CeilIntToInt:
		f = &builtinCeilIntToIntSig{base}
	case tipb.ScalarFuncSig_CeilDecToInt:
		f = &builtinCeilDecToIntSig{base}
	case tipb.ScalarFuncSig_CeilDecToDec:
		f = &builtinCeilDecToDecSig{base}
	case tipb.ScalarFuncSig_CeilReal:
		f = &builtinCeilRealSig{base}
	case tipb.ScalarFuncSig_FloorIntToDec:
		f = &builtinFloorIntToDecSig{base}
	case tipb.ScalarFuncSig_FloorIntToInt:
		f = &builtinFloorIntToIntSig{base}
	case tipb.ScalarFuncSig_FloorDecToInt:
		f = &builtinFloorDecToIntSig{base}
	case tipb.ScalarFuncSig_FloorDecToDec:
		f = &builtinFloorDecToDecSig{base}
	case tipb.ScalarFuncSig_FloorReal:
		f = &builtinFloorRealSig{base}
	case tipb.ScalarFuncSig_RoundReal:
		f = &builtinRoundRealSig{base}
	case tipb.ScalarFuncSig_RoundInt:
		f = &builtinRoundIntSig{base}
	case tipb.ScalarFuncSig_RoundDec:
		f = &builtinRoundDecSig{base}
	case tipb.ScalarFuncSig_RoundWithFracReal:
		f = &builtinRoundWithFracRealSig{base}
	case tipb.ScalarFuncSig_RoundWithFracInt:
		f = &builtinRoundWithFracIntSig{base}
	case tipb.ScalarFuncSig_RoundWithFracDec:
		f = &builtinRoundWithFracDecSig{base}
	case tipb.ScalarFuncSig_Log1Arg:
		f = &builtinLog1ArgSig{base}
	case tipb.ScalarFuncSig_Log2Args:
		f = &builtinLog2ArgsSig{base}
	case tipb.ScalarFuncSig_Log2:
		f = &builtinLog2Sig{base}
	case tipb.ScalarFuncSig_Log10:
		f = &builtinLog10Sig{base}
	// case tipb.ScalarFuncSig_Rand:
	case tipb.ScalarFuncSig_RandWithSeedFirstGen:
		f = &builtinRandWithSeedFirstGenSig{base}
	case tipb.ScalarFuncSig_Pow:
		f = &builtinPowSig{base}
	case tipb.ScalarFuncSig_Conv:
		f = &builtinConvSig{base}
	case tipb.ScalarFuncSig_CRC32:
		f = &builtinCRC32Sig{base}
	case tipb.ScalarFuncSig_Sign:
		f = &builtinSignSig{base}
	case tipb.ScalarFuncSig_Sqrt:
		f = &builtinSqrtSig{base}
	case tipb.ScalarFuncSig_Acos:
		f = &builtinAcosSig{base}
	case tipb.ScalarFuncSig_Asin:
		f = &builtinAsinSig{base}
	case tipb.ScalarFuncSig_Atan1Arg:
		f = &builtinAtan1ArgSig{base}
	case tipb.ScalarFuncSig_Atan2Args:
		f = &builtinAtan2ArgsSig{base}
	case tipb.ScalarFuncSig_Cos:
		f = &builtinCosSig{base}
	case tipb.ScalarFuncSig_Cot:
		f = &builtinCotSig{base}
	case tipb.ScalarFuncSig_Degrees:
		f = &builtinDegreesSig{base}
	case tipb.ScalarFuncSig_Exp:
		f = &builtinExpSig{base}
	case tipb.ScalarFuncSig_PI:
		f = &builtinPISig{base}
	case tipb.ScalarFuncSig_Radians:
		f = &builtinRadiansSig{base}
	case tipb.ScalarFuncSig_Sin:
		f = &builtinSinSig{base}
	case tipb.ScalarFuncSig_Tan:
		f = &builtinTanSig{base}
	case tipb.ScalarFuncSig_TruncateInt:
		f = &builtinTruncateIntSig{base}
	case tipb.ScalarFuncSig_TruncateReal:
		f = &builtinTruncateRealSig{base}
	case tipb.ScalarFuncSig_TruncateDecimal:
		f = &builtinTruncateDecimalSig{base}
	case tipb.ScalarFuncSig_TruncateUint:
		f = &builtinTruncateUintSig{base}
	case tipb.ScalarFuncSig_LogicalAnd:
		f = &builtinLogicAndSig{base}
	case tipb.ScalarFuncSig_LogicalOr:
		f = &builtinLogicOrSig{base}
	case tipb.ScalarFuncSig_LogicalXor:
		f = &builtinLogicXorSig{base}
	case tipb.ScalarFuncSig_UnaryNotInt:
		f = &builtinUnaryNotIntSig{base}
	case tipb.ScalarFuncSig_UnaryNotDecimal:
		f = &builtinUnaryNotDecimalSig{base}
	case tipb.ScalarFuncSig_UnaryNotReal:
		f = &builtinUnaryNotRealSig{base}
	case tipb.ScalarFuncSig_UnaryMinusInt:
		f = &builtinUnaryMinusIntSig{base}
	case tipb.ScalarFuncSig_UnaryMinusReal:
		f = &builtinUnaryMinusRealSig{base}
	case tipb.ScalarFuncSig_UnaryMinusDecimal:
		f = &builtinUnaryMinusDecimalSig{base, false}
	case tipb.ScalarFuncSig_DecimalIsNull:
		f = &builtinDecimalIsNullSig{base}
	case tipb.ScalarFuncSig_DurationIsNull:
		f = &builtinDurationIsNullSig{base}
	case tipb.ScalarFuncSig_RealIsNull:
		f = &builtinRealIsNullSig{base}
	case tipb.ScalarFuncSig_StringIsNull:
		f = &builtinStringIsNullSig{base}
	case tipb.ScalarFuncSig_TimeIsNull:
		f = &builtinTimeIsNullSig{base}
	case tipb.ScalarFuncSig_IntIsNull:
		f = &builtinIntIsNullSig{base}
	// case tipb.ScalarFuncSig_JsonIsNull:
	case tipb.ScalarFuncSig_BitAndSig:
		f = &builtinBitAndSig{base}
	case tipb.ScalarFuncSig_BitOrSig:
		f = &builtinBitOrSig{base}
	case tipb.ScalarFuncSig_BitXorSig:
		f = &builtinBitXorSig{base}
	case tipb.ScalarFuncSig_BitNegSig:
		f = &builtinBitNegSig{base}
	case tipb.ScalarFuncSig_IntIsTrue:
		f = &builtinIntIsTrueSig{base, false}
	case tipb.ScalarFuncSig_RealIsTrue:
		f = &builtinRealIsTrueSig{base, false}
	case tipb.ScalarFuncSig_DecimalIsTrue:
		f = &builtinDecimalIsTrueSig{base, false}
	case tipb.ScalarFuncSig_IntIsFalse:
		f = &builtinIntIsFalseSig{base, false}
	case tipb.ScalarFuncSig_RealIsFalse:
		f = &builtinRealIsFalseSig{base, false}
	case tipb.ScalarFuncSig_DecimalIsFalse:
		f = &builtinDecimalIsFalseSig{base, false}
	case tipb.ScalarFuncSig_IntIsTrueWithNull:
		f = &builtinIntIsTrueSig{base, true}
	case tipb.ScalarFuncSig_RealIsTrueWithNull:
		f = &builtinRealIsTrueSig{base, true}
	case tipb.ScalarFuncSig_DecimalIsTrueWithNull:
		f = &builtinDecimalIsTrueSig{base, true}
	case tipb.ScalarFuncSig_IntIsFalseWithNull:
		f = &builtinIntIsFalseSig{base, true}
	case tipb.ScalarFuncSig_RealIsFalseWithNull:
		f = &builtinRealIsFalseSig{base, true}
	case tipb.ScalarFuncSig_DecimalIsFalseWithNull:
		f = &builtinDecimalIsFalseSig{base, true}
	case tipb.ScalarFuncSig_LeftShift:
		f = &builtinLeftShiftSig{base}
	case tipb.ScalarFuncSig_RightShift:
		f = &builtinRightShiftSig{base}
	case tipb.ScalarFuncSig_BitCount:
		f = &builtinBitCountSig{base}
	case tipb.ScalarFuncSig_GetParamString:
		f = &builtinGetParamStringSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_GetVar:
		f = &builtinGetStringVarSig{baseBuiltinFunc: base}
	// case tipb.ScalarFuncSig_RowSig:
	case tipb.ScalarFuncSig_SetVar:
		f = &builtinSetStringVarSig{baseBuiltinFunc: base}
	// case tipb.ScalarFuncSig_ValuesDecimal:
	// 	f = &builtinValuesDecimalSig{base}
	// case tipb.ScalarFuncSig_ValuesDuration:
	// 	f = &builtinValuesDurationSig{base}
	// case tipb.ScalarFuncSig_ValuesInt:
	// 	f = &builtinValuesIntSig{base}
	// case tipb.ScalarFuncSig_ValuesJSON:
	// 	f = &builtinValuesJSONSig{base}
	// case tipb.ScalarFuncSig_ValuesReal:
	// 	f = &builtinValuesRealSig{base}
	// case tipb.ScalarFuncSig_ValuesString:
	// 	f = &builtinValuesStringSig{base}
	// case tipb.ScalarFuncSig_ValuesTime:
	// 	f = &builtinValuesTimeSig{base}
	case tipb.ScalarFuncSig_InInt:
		f = &builtinInIntSig{baseInSig: baseInSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_InReal:
		f = &builtinInRealSig{baseInSig: baseInSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_InDecimal:
		f = &builtinInDecimalSig{baseInSig: baseInSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_InString:
		f = &builtinInStringSig{baseInSig: baseInSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_InTime:
		f = &builtinInTimeSig{baseInSig: baseInSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_InDuration:
		f = &builtinInDurationSig{baseInSig: baseInSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_InJson:
		f = &builtinInJSONSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_IfNullInt:
		f = &builtinIfNullIntSig{base}
	case tipb.ScalarFuncSig_IfNullReal:
		f = &builtinIfNullRealSig{base}
	case tipb.ScalarFuncSig_IfNullDecimal:
		f = &builtinIfNullDecimalSig{base}
	case tipb.ScalarFuncSig_IfNullString:
		f = &builtinIfNullStringSig{base}
	case tipb.ScalarFuncSig_IfNullTime:
		f = &builtinIfNullTimeSig{base}
	case tipb.ScalarFuncSig_IfNullDuration:
		f = &builtinIfNullDurationSig{base}
	case tipb.ScalarFuncSig_IfInt:
		f = &builtinIfIntSig{base}
	case tipb.ScalarFuncSig_IfReal:
		f = &builtinIfRealSig{base}
	case tipb.ScalarFuncSig_IfDecimal:
		f = &builtinIfDecimalSig{base}
	case tipb.ScalarFuncSig_IfString:
		f = &builtinIfStringSig{base}
	case tipb.ScalarFuncSig_IfTime:
		f = &builtinIfTimeSig{base}
	case tipb.ScalarFuncSig_IfDuration:
		f = &builtinIfDurationSig{base}
	case tipb.ScalarFuncSig_IfNullJson:
		f = &builtinIfNullJSONSig{base}
	case tipb.ScalarFuncSig_IfJson:
		f = &builtinIfJSONSig{base}
	case tipb.ScalarFuncSig_CaseWhenInt:
		f = &builtinCaseWhenIntSig{base}
	case tipb.ScalarFuncSig_CaseWhenReal:
		f = &builtinCaseWhenRealSig{base}
	case tipb.ScalarFuncSig_CaseWhenDecimal:
		f = &builtinCaseWhenDecimalSig{base}
	case tipb.ScalarFuncSig_CaseWhenString:
		f = &builtinCaseWhenStringSig{base}
	case tipb.ScalarFuncSig_CaseWhenTime:
		f = &builtinCaseWhenTimeSig{base}
	case tipb.ScalarFuncSig_CaseWhenDuration:
		f = &builtinCaseWhenDurationSig{base}
	case tipb.ScalarFuncSig_CaseWhenJson:
		f = &builtinCaseWhenJSONSig{base}
	// case tipb.ScalarFuncSig_AesDecrypt:
	// 	f = &builtinAesDecryptSig{base}
	// case tipb.ScalarFuncSig_AesEncrypt:
	// 	f = &builtinAesEncryptSig{base}
	case tipb.ScalarFuncSig_Compress:
		f = &builtinCompressSig{base}
	case tipb.ScalarFuncSig_MD5:
		f = &builtinMD5Sig{base}
	case tipb.ScalarFuncSig_Password:
		f = &builtinPasswordSig{base}
	case tipb.ScalarFuncSig_RandomBytes:
		f = &builtinRandomBytesSig{base}
	case tipb.ScalarFuncSig_SHA1:
		f = &builtinSHA1Sig{base}
	case tipb.ScalarFuncSig_SHA2:
		f = &builtinSHA2Sig{base}
	case tipb.ScalarFuncSig_Uncompress:
		f = &builtinUncompressSig{base}
	case tipb.ScalarFuncSig_UncompressedLength:
		f = &builtinUncompressedLengthSig{base}
	case tipb.ScalarFuncSig_Database:
		f = &builtinDatabaseSig{base}
	case tipb.ScalarFuncSig_FoundRows:
		f = &builtinFoundRowsSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_CurrentUser:
		f = &builtinCurrentUserSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_User:
		f = &builtinUserSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_ConnectionID:
		f = &builtinConnectionIDSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_LastInsertID:
		f = &builtinLastInsertIDSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_LastInsertIDWithID:
		f = &builtinLastInsertIDWithIDSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_Version:
		f = &builtinVersionSig{base}
	case tipb.ScalarFuncSig_TiDBVersion:
		f = &builtinTiDBVersionSig{base}
	case tipb.ScalarFuncSig_RowCount:
		f = &builtinRowCountSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_Sleep:
		f = &builtinSleepSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_Lock:
		f = &builtinLockSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_ReleaseLock:
		f = &builtinReleaseLockSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_DecimalAnyValue:
		f = &builtinDecimalAnyValueSig{base}
	case tipb.ScalarFuncSig_DurationAnyValue:
		f = &builtinDurationAnyValueSig{base}
	case tipb.ScalarFuncSig_IntAnyValue:
		f = &builtinIntAnyValueSig{base}
	case tipb.ScalarFuncSig_JSONAnyValue:
		f = &builtinJSONAnyValueSig{base}
	case tipb.ScalarFuncSig_RealAnyValue:
		f = &builtinRealAnyValueSig{base}
	case tipb.ScalarFuncSig_StringAnyValue:
		f = &builtinStringAnyValueSig{base}
	case tipb.ScalarFuncSig_TimeAnyValue:
		f = &builtinTimeAnyValueSig{base}
	case tipb.ScalarFuncSig_InetAton:
		f = &builtinInetAtonSig{base}
	case tipb.ScalarFuncSig_InetNtoa:
		f = &builtinInetNtoaSig{base}
	case tipb.ScalarFuncSig_Inet6Aton:
		f = &builtinInet6AtonSig{base}
	case tipb.ScalarFuncSig_Inet6Ntoa:
		f = &builtinInet6NtoaSig{base}
	case tipb.ScalarFuncSig_IsIPv4:
		f = &builtinIsIPv4Sig{base}
	case tipb.ScalarFuncSig_IsIPv4Compat:
		f = &builtinIsIPv4CompatSig{base}
	case tipb.ScalarFuncSig_IsIPv4Mapped:
		f = &builtinIsIPv4MappedSig{base}
	case tipb.ScalarFuncSig_IsIPv6:
		f = &builtinIsIPv6Sig{base}
	case tipb.ScalarFuncSig_UUID:
		f = &builtinUUIDSig{base}
	case tipb.ScalarFuncSig_LikeSig:
		f = &builtinLikeSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_IlikeSig:
		f = &builtinIlikeSig{baseBuiltinFunc: base}
	case tipb.ScalarFuncSig_RegexpSig:
		f = &builtinRegexpLikeFuncSig{regexpBaseFuncSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_RegexpUTF8Sig:
		f = &builtinRegexpLikeFuncSig{regexpBaseFuncSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_RegexpLikeSig:
		f = &builtinRegexpLikeFuncSig{regexpBaseFuncSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_RegexpSubstrSig:
		f = &builtinRegexpSubstrFuncSig{regexpBaseFuncSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_RegexpInStrSig:
		f = &builtinRegexpInStrFuncSig{regexpBaseFuncSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_RegexpReplaceSig:
		f = &builtinRegexpReplaceFuncSig{regexpBaseFuncSig: regexpBaseFuncSig{baseBuiltinFunc: base}}
	case tipb.ScalarFuncSig_JsonExtractSig:
		f = &builtinJSONExtractSig{base}
	case tipb.ScalarFuncSig_JsonUnquoteSig:
		f = &builtinJSONUnquoteSig{base}
	case tipb.ScalarFuncSig_JsonTypeSig:
		f = &builtinJSONTypeSig{base}
	case tipb.ScalarFuncSig_JsonSetSig:
		f = &builtinJSONSetSig{base}
	case tipb.ScalarFuncSig_JsonInsertSig:
		f = &builtinJSONInsertSig{base}
	case tipb.ScalarFuncSig_JsonReplaceSig:
		f = &builtinJSONReplaceSig{base}
	case tipb.ScalarFuncSig_JsonRemoveSig:
		f = &builtinJSONRemoveSig{base}
	case tipb.ScalarFuncSig_JsonMergeSig:
		f = &builtinJSONMergeSig{base}
	case tipb.ScalarFuncSig_JsonObjectSig:
		f = &builtinJSONObjectSig{base}
	case tipb.ScalarFuncSig_JsonArraySig:
		f = &builtinJSONArraySig{base}
	case tipb.ScalarFuncSig_JsonValidJsonSig:
		f = &builtinJSONValidJSONSig{base}
	case tipb.ScalarFuncSig_JsonContainsSig:
		f = &builtinJSONContainsSig{base}
	case tipb.ScalarFuncSig_JsonArrayAppendSig:
		f = &builtinJSONArrayAppendSig{base}
	case tipb.ScalarFuncSig_JsonArrayInsertSig:
		f = &builtinJSONArrayInsertSig{base}
	case tipb.ScalarFuncSig_JsonMergePatchSig:
		f = &builtinJSONMergePatchSig{base}
	case tipb.ScalarFuncSig_JsonMergePreserveSig:
		f = &builtinJSONMergeSig{base}
	case tipb.ScalarFuncSig_JsonContainsPathSig:
		f = &builtinJSONContainsPathSig{base}
	// case tipb.ScalarFuncSig_JsonPrettySig:
	case tipb.ScalarFuncSig_JsonQuoteSig:
		f = &builtinJSONQuoteSig{base}
	case tipb.ScalarFuncSig_JsonSearchSig:
		f = &builtinJSONSearchSig{base}
	case tipb.ScalarFuncSig_JsonStorageSizeSig:
		f = &builtinJSONStorageSizeSig{base}
	case tipb.ScalarFuncSig_JsonDepthSig:
		f = &builtinJSONDepthSig{base}
	case tipb.ScalarFuncSig_JsonKeysSig:
		f = &builtinJSONKeysSig{base}
	case tipb.ScalarFuncSig_JsonLengthSig:
		f = &builtinJSONLengthSig{base}
	case tipb.ScalarFuncSig_JsonKeys2ArgsSig:
		f = &builtinJSONKeys2ArgsSig{base}
	case tipb.ScalarFuncSig_JsonValidStringSig:
		f = &builtinJSONValidStringSig{base}
	case tipb.ScalarFuncSig_JsonValidOthersSig:
		f = &builtinJSONValidOthersSig{base}
	case tipb.ScalarFuncSig_JsonMemberOfSig:
		f = &builtinJSONMemberOfSig{base}
	case tipb.ScalarFuncSig_DateFormatSig:
		f = &builtinDateFormatSig{base}
	// case tipb.ScalarFuncSig_DateLiteral:
	// 	f = &builtinDateLiteralSig{base}
	case tipb.ScalarFuncSig_DateDiff:
		f = &builtinDateDiffSig{base}
	case tipb.ScalarFuncSig_NullTimeDiff:
		f = &builtinNullTimeDiffSig{base}
	case tipb.ScalarFuncSig_TimeStringTimeDiff:
		f = &builtinTimeStringTimeDiffSig{base}
	case tipb.ScalarFuncSig_DurationStringTimeDiff:
		f = &builtinDurationStringTimeDiffSig{base}
	case tipb.ScalarFuncSig_DurationDurationTimeDiff:
		f = &builtinDurationDurationTimeDiffSig{base}
	case tipb.ScalarFuncSig_StringTimeTimeDiff:
		f = &builtinStringTimeTimeDiffSig{base}
	case tipb.ScalarFuncSig_StringDurationTimeDiff:
		f = &builtinStringDurationTimeDiffSig{base}
	case tipb.ScalarFuncSig_StringStringTimeDiff:
		f = &builtinStringStringTimeDiffSig{base}
	case tipb.ScalarFuncSig_TimeTimeTimeDiff:
		f = &builtinTimeTimeTimeDiffSig{base}
	case tipb.ScalarFuncSig_Date:
		f = &builtinDateSig{base}
	case tipb.ScalarFuncSig_Hour:
		f = &builtinHourSig{base}
	case tipb.ScalarFuncSig_Minute:
		f = &builtinMinuteSig{base}
	case tipb.ScalarFuncSig_Second:
		f = &builtinSecondSig{base}
	case tipb.ScalarFuncSig_MicroSecond:
		f = &builtinMicroSecondSig{base}
	case tipb.ScalarFuncSig_Month:
		f = &builtinMonthSig{base}
	case tipb.ScalarFuncSig_MonthName:
		f = &builtinMonthNameSig{base}
	case tipb.ScalarFuncSig_NowWithArg:
		f = &builtinNowWithArgSig{base}
	case tipb.ScalarFuncSig_NowWithoutArg:
		f = &builtinNowWithoutArgSig{base}
	case tipb.ScalarFuncSig_DayName:
		f = &builtinDayNameSig{base}
	case tipb.ScalarFuncSig_DayOfMonth:
		f = &builtinDayOfMonthSig{base}
	case tipb.ScalarFuncSig_DayOfWeek:
		f = &builtinDayOfWeekSig{base}
	case tipb.ScalarFuncSig_DayOfYear:
		f = &builtinDayOfYearSig{base}
	case tipb.ScalarFuncSig_WeekWithMode:
		f = &builtinWeekWithModeSig{base}
	case tipb.ScalarFuncSig_WeekWithoutMode:
		f = &builtinWeekWithoutModeSig{base}
	case tipb.ScalarFuncSig_WeekDay:
		f = &builtinWeekDaySig{base}
	case tipb.ScalarFuncSig_WeekOfYear:
		f = &builtinWeekOfYearSig{base}
	case tipb.ScalarFuncSig_Year:
		f = &builtinYearSig{base}
	case tipb.ScalarFuncSig_YearWeekWithMode:
		f = &builtinYearWeekWithModeSig{base}
	case tipb.ScalarFuncSig_YearWeekWithoutMode:
		f = &builtinYearWeekWithoutModeSig{base}
	case tipb.ScalarFuncSig_GetFormat:
		f = &builtinGetFormatSig{base}
	case tipb.ScalarFuncSig_SysDateWithFsp:
		f = &builtinSysDateWithFspSig{base}
	case tipb.ScalarFuncSig_SysDateWithoutFsp:
		f = &builtinSysDateWithoutFspSig{base}
	case tipb.ScalarFuncSig_CurrentDate:
		f = &builtinCurrentDateSig{base}
	case tipb.ScalarFuncSig_CurrentTime0Arg:
		f = &builtinCurrentTime0ArgSig{base}
	case tipb.ScalarFuncSig_CurrentTime1Arg:
		f = &builtinCurrentTime1ArgSig{base}
	case tipb.ScalarFuncSig_Time:
		f = &builtinTimeSig{base}
	// case tipb.ScalarFuncSig_TimeLiteral:
	// 	f = &builtinTimeLiteralSig{base}
	case tipb.ScalarFuncSig_UTCDate:
		f = &builtinUTCDateSig{base}
	case tipb.ScalarFuncSig_UTCTimestampWithArg:
		f = &builtinUTCTimestampWithArgSig{base}
	case tipb.ScalarFuncSig_UTCTimestampWithoutArg:
		f = &builtinUTCTimestampWithoutArgSig{base}
	case tipb.ScalarFuncSig_AddDatetimeAndDuration:
		f = &builtinAddDatetimeAndDurationSig{base}
	case tipb.ScalarFuncSig_AddDatetimeAndString:
		f = &builtinAddDatetimeAndStringSig{base}
	case tipb.ScalarFuncSig_AddTimeDateTimeNull:
		f = &builtinAddTimeDateTimeNullSig{base}
	case tipb.ScalarFuncSig_AddStringAndDuration:
		f = &builtinAddStringAndDurationSig{base}
	case tipb.ScalarFuncSig_AddStringAndString:
		f = &builtinAddStringAndStringSig{base}
	case tipb.ScalarFuncSig_AddTimeStringNull:
		f = &builtinAddTimeStringNullSig{base}
	case tipb.ScalarFuncSig_AddDurationAndDuration:
		f = &builtinAddDurationAndDurationSig{base}
	case tipb.ScalarFuncSig_AddDurationAndString:
		f = &builtinAddDurationAndStringSig{base}
	case tipb.ScalarFuncSig_AddTimeDurationNull:
		f = &builtinAddTimeDurationNullSig{base}
	case tipb.ScalarFuncSig_AddDateAndDuration:
		f = &builtinAddDateAndDurationSig{base}
	case tipb.ScalarFuncSig_AddDateAndString:
		f = &builtinAddDateAndStringSig{base}
	case tipb.ScalarFuncSig_SubDatetimeAndDuration:
		f = &builtinSubDatetimeAndDurationSig{base}
	case tipb.ScalarFuncSig_SubDatetimeAndString:
		f = &builtinSubDatetimeAndStringSig{base}
	case tipb.ScalarFuncSig_SubTimeDateTimeNull:
		f = &builtinSubTimeDateTimeNullSig{base}
	case tipb.ScalarFuncSig_SubStringAndDuration:
		f = &builtinSubStringAndDurationSig{base}
	case tipb.ScalarFuncSig_SubStringAndString:
		f = &builtinSubStringAndStringSig{base}
	case tipb.ScalarFuncSig_SubTimeStringNull:
		f = &builtinSubTimeStringNullSig{base}
	case tipb.ScalarFuncSig_SubDurationAndDuration:
		f = &builtinSubDurationAndDurationSig{base}
	case tipb.ScalarFuncSig_SubDurationAndString:
		f = &builtinSubDurationAndStringSig{base}
	case tipb.ScalarFuncSig_SubTimeDurationNull:
		f = &builtinSubTimeDurationNullSig{base}
	case tipb.ScalarFuncSig_SubDateAndDuration:
		f = &builtinSubDateAndDurationSig{base}
	case tipb.ScalarFuncSig_SubDateAndString:
		f = &builtinSubDateAndStringSig{base}
	case tipb.ScalarFuncSig_UnixTimestampCurrent:
		f = &builtinUnixTimestampCurrentSig{base}
	case tipb.ScalarFuncSig_UnixTimestampInt:
		f = &builtinUnixTimestampIntSig{base}
	case tipb.ScalarFuncSig_UnixTimestampDec:
		f = &builtinUnixTimestampDecSig{base}
	// case tipb.ScalarFuncSig_ConvertTz:
	// 	f = &builtinConvertTzSig{base}
	case tipb.ScalarFuncSig_MakeDate:
		f = &builtinMakeDateSig{base}
	case tipb.ScalarFuncSig_MakeTime:
		f = &builtinMakeTimeSig{base}
	case tipb.ScalarFuncSig_PeriodAdd:
		f = &builtinPeriodAddSig{base}
	case tipb.ScalarFuncSig_PeriodDiff:
		f = &builtinPeriodDiffSig{base}
	case tipb.ScalarFuncSig_Quarter:
		f = &builtinQuarterSig{base}
	case tipb.ScalarFuncSig_SecToTime:
		f = &builtinSecToTimeSig{base}
	case tipb.ScalarFuncSig_TimeToSec:
		f = &builtinTimeToSecSig{base}
	case tipb.ScalarFuncSig_TimestampAdd:
		f = &builtinTimestampAddSig{base}
	case tipb.ScalarFuncSig_ToDays:
		f = &builtinToDaysSig{base}
	case tipb.ScalarFuncSig_ToSeconds:
		f = &builtinToSecondsSig{base}
	case tipb.ScalarFuncSig_UTCTimeWithArg:
		f = &builtinUTCTimeWithArgSig{base}
	case tipb.ScalarFuncSig_UTCTimeWithoutArg:
		f = &builtinUTCTimeWithoutArgSig{base}
	// case tipb.ScalarFuncSig_Timestamp1Arg:
	// 	f = &builtinTimestamp1ArgSig{base}
	// case tipb.ScalarFuncSig_Timestamp2Args:
	// 	f = &builtinTimestamp2ArgsSig{base}
	// case tipb.ScalarFuncSig_TimestampLiteral:
	// 	f = &builtinTimestampLiteralSig{base}
	case tipb.ScalarFuncSig_LastDay:
		f = &builtinLastDaySig{base}
	case tipb.ScalarFuncSig_StrToDateDate:
		f = &builtinStrToDateDateSig{base}
	case tipb.ScalarFuncSig_StrToDateDatetime:
		f = &builtinStrToDateDatetimeSig{base}
	case tipb.ScalarFuncSig_StrToDateDuration:
		f = &builtinStrToDateDurationSig{base}
	case tipb.ScalarFuncSig_FromUnixTime1Arg:
		f = &builtinFromUnixTime1ArgSig{base}
	case tipb.ScalarFuncSig_FromUnixTime2Arg:
		f = &builtinFromUnixTime2ArgSig{base}
	case tipb.ScalarFuncSig_ExtractDatetimeFromString:
		f = &builtinExtractDatetimeFromStringSig{base}
	case tipb.ScalarFuncSig_ExtractDatetime:
		f = &builtinExtractDatetimeSig{base}
	case tipb.ScalarFuncSig_ExtractDuration:
		f = &builtinExtractDurationSig{base}
	case tipb.ScalarFuncSig_AddDateStringString,
		tipb.ScalarFuncSig_AddDateStringInt,
		tipb.ScalarFuncSig_AddDateStringReal,
		tipb.ScalarFuncSig_AddDateStringDecimal,
		tipb.ScalarFuncSig_AddDateIntString,
		tipb.ScalarFuncSig_AddDateIntInt,
		tipb.ScalarFuncSig_AddDateIntReal,
		tipb.ScalarFuncSig_AddDateIntDecimal,
		tipb.ScalarFuncSig_AddDateRealString,
		tipb.ScalarFuncSig_AddDateRealInt,
		tipb.ScalarFuncSig_AddDateRealReal,
		tipb.ScalarFuncSig_AddDateRealDecimal,
		tipb.ScalarFuncSig_AddDateDecimalString,
		tipb.ScalarFuncSig_AddDateDecimalInt,
		tipb.ScalarFuncSig_AddDateDecimalReal,
		tipb.ScalarFuncSig_AddDateDecimalDecimal,
		tipb.ScalarFuncSig_AddDateDatetimeString,
		tipb.ScalarFuncSig_AddDateDatetimeInt,
		tipb.ScalarFuncSig_AddDateDatetimeReal,
		tipb.ScalarFuncSig_AddDateDatetimeDecimal,
		tipb.ScalarFuncSig_AddDateDurationString,
		tipb.ScalarFuncSig_AddDateDurationInt,
		tipb.ScalarFuncSig_AddDateDurationReal,
		tipb.ScalarFuncSig_AddDateDurationDecimal,
		tipb.ScalarFuncSig_AddDateDurationStringDatetime,
		tipb.ScalarFuncSig_AddDateDurationIntDatetime,
		tipb.ScalarFuncSig_AddDateDurationRealDatetime,
		tipb.ScalarFuncSig_AddDateDurationDecimalDatetime:
		f, e = (&addSubDateFunctionClass{baseFunctionClass{ast.AddDate, 3, 3}, addTime, addDuration, setAdd}).getFunction(ctx, args)
		if e != nil {
			return f, e
		}
	case tipb.ScalarFuncSig_SubDateStringString,
		tipb.ScalarFuncSig_SubDateStringInt,
		tipb.ScalarFuncSig_SubDateStringReal,
		tipb.ScalarFuncSig_SubDateStringDecimal,
		tipb.ScalarFuncSig_SubDateIntString,
		tipb.ScalarFuncSig_SubDateIntInt,
		tipb.ScalarFuncSig_SubDateIntReal,
		tipb.ScalarFuncSig_SubDateIntDecimal,
		tipb.ScalarFuncSig_SubDateRealString,
		tipb.ScalarFuncSig_SubDateRealInt,
		tipb.ScalarFuncSig_SubDateRealReal,
		tipb.ScalarFuncSig_SubDateRealDecimal,
		tipb.ScalarFuncSig_SubDateDecimalString,
		tipb.ScalarFuncSig_SubDateDecimalInt,
		tipb.ScalarFuncSig_SubDateDecimalReal,
		tipb.ScalarFuncSig_SubDateDecimalDecimal,
		tipb.ScalarFuncSig_SubDateDatetimeString,
		tipb.ScalarFuncSig_SubDateDatetimeInt,
		tipb.ScalarFuncSig_SubDateDatetimeReal,
		tipb.ScalarFuncSig_SubDateDatetimeDecimal,
		tipb.ScalarFuncSig_SubDateDurationString,
		tipb.ScalarFuncSig_SubDateDurationInt,
		tipb.ScalarFuncSig_SubDateDurationReal,
		tipb.ScalarFuncSig_SubDateDurationDecimal,
		tipb.ScalarFuncSig_SubDateDurationStringDatetime,
		tipb.ScalarFuncSig_SubDateDurationIntDatetime,
		tipb.ScalarFuncSig_SubDateDurationRealDatetime,
		tipb.ScalarFuncSig_SubDateDurationDecimalDatetime:
		f, e = (&addSubDateFunctionClass{baseFunctionClass{ast.SubDate, 3, 3}, subTime, subDuration, setSub}).getFunction(ctx, args)
		if e != nil {
			return f, e
		}
	case tipb.ScalarFuncSig_FromDays:
		f = &builtinFromDaysSig{base}
	case tipb.ScalarFuncSig_TimeFormat:
		f = &builtinTimeFormatSig{base}
	case tipb.ScalarFuncSig_TimestampDiff:
		f = &builtinTimestampDiffSig{base}
	case tipb.ScalarFuncSig_BitLength:
		f = &builtinBitLengthSig{base}
	case tipb.ScalarFuncSig_Bin:
		f = &builtinBinSig{base}
	case tipb.ScalarFuncSig_ASCII:
		f = &builtinASCIISig{base}
	case tipb.ScalarFuncSig_Char:
		f = &builtinCharSig{base}
	case tipb.ScalarFuncSig_CharLengthUTF8:
		f = &builtinCharLengthUTF8Sig{base}
	case tipb.ScalarFuncSig_CharLength:
		f = &builtinCharLengthBinarySig{base}
	case tipb.ScalarFuncSig_Concat:
		f = &builtinConcatSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_ConcatWS:
		f = &builtinConcatWSSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Convert:
		f = &builtinConvertSig{base}
	case tipb.ScalarFuncSig_Elt:
		f = &builtinEltSig{base}
	case tipb.ScalarFuncSig_ExportSet3Arg:
		f = &builtinExportSet3ArgSig{base}
	case tipb.ScalarFuncSig_ExportSet4Arg:
		f = &builtinExportSet4ArgSig{base}
	case tipb.ScalarFuncSig_ExportSet5Arg:
		f = &builtinExportSet5ArgSig{base}
	case tipb.ScalarFuncSig_FieldInt:
		f = &builtinFieldIntSig{base}
	case tipb.ScalarFuncSig_FieldReal:
		f = &builtinFieldRealSig{base}
	case tipb.ScalarFuncSig_FieldString:
		f = &builtinFieldStringSig{base}
	case tipb.ScalarFuncSig_FindInSet:
		f = &builtinFindInSetSig{base}
	case tipb.ScalarFuncSig_Format:
		f = &builtinFormatSig{base}
	case tipb.ScalarFuncSig_FormatWithLocale:
		f = &builtinFormatWithLocaleSig{base}
	case tipb.ScalarFuncSig_FromBase64:
		f = &builtinFromBase64Sig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_HexIntArg:
		f = &builtinHexIntArgSig{base}
	case tipb.ScalarFuncSig_HexStrArg:
		f = &builtinHexStrArgSig{base}
	case tipb.ScalarFuncSig_InsertUTF8:
		f = &builtinInsertUTF8Sig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Insert:
		f = &builtinInsertSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_InstrUTF8:
		f = &builtinInstrUTF8Sig{base}
	case tipb.ScalarFuncSig_Instr:
		f = &builtinInstrSig{base}
	case tipb.ScalarFuncSig_LTrim:
		f = &builtinLTrimSig{base}
	case tipb.ScalarFuncSig_LeftUTF8:
		f = &builtinLeftUTF8Sig{base}
	case tipb.ScalarFuncSig_Left:
		f = &builtinLeftSig{base}
	case tipb.ScalarFuncSig_Length:
		f = &builtinLengthSig{base}
	case tipb.ScalarFuncSig_Locate2ArgsUTF8:
		f = &builtinLocate2ArgsUTF8Sig{base}
	case tipb.ScalarFuncSig_Locate3ArgsUTF8:
		f = &builtinLocate3ArgsUTF8Sig{base}
	case tipb.ScalarFuncSig_Locate2Args:
		f = &builtinLocate2ArgsSig{base}
	case tipb.ScalarFuncSig_Locate3Args:
		f = &builtinLocate3ArgsSig{base}
	case tipb.ScalarFuncSig_Lower:
		f = &builtinLowerSig{base}
	case tipb.ScalarFuncSig_LowerUTF8:
		f = &builtinLowerUTF8Sig{base}
	case tipb.ScalarFuncSig_LpadUTF8:
		f = &builtinLpadUTF8Sig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Lpad:
		f = &builtinLpadSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_MakeSet:
		f = &builtinMakeSetSig{base}
	case tipb.ScalarFuncSig_OctInt:
		f = &builtinOctIntSig{base}
	case tipb.ScalarFuncSig_OctString:
		f = &builtinOctStringSig{base}
	case tipb.ScalarFuncSig_Ord:
		f = &builtinOrdSig{base}
	case tipb.ScalarFuncSig_Quote:
		f = &builtinQuoteSig{base}
	case tipb.ScalarFuncSig_RTrim:
		f = &builtinRTrimSig{base}
	case tipb.ScalarFuncSig_Repeat:
		f = &builtinRepeatSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Replace:
		f = &builtinReplaceSig{base}
	case tipb.ScalarFuncSig_ReverseUTF8:
		f = &builtinReverseUTF8Sig{base}
	case tipb.ScalarFuncSig_Reverse:
		f = &builtinReverseSig{base}
	case tipb.ScalarFuncSig_RightUTF8:
		f = &builtinRightUTF8Sig{base}
	case tipb.ScalarFuncSig_Right:
		f = &builtinRightSig{base}
	case tipb.ScalarFuncSig_RpadUTF8:
		f = &builtinRpadUTF8Sig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Rpad:
		f = &builtinRpadSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Space:
		f = &builtinSpaceSig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Strcmp:
		f = &builtinStrcmpSig{base}
	case tipb.ScalarFuncSig_Substring2ArgsUTF8:
		f = &builtinSubstring2ArgsUTF8Sig{base}
	case tipb.ScalarFuncSig_Substring3ArgsUTF8:
		f = &builtinSubstring3ArgsUTF8Sig{base}
	case tipb.ScalarFuncSig_Substring2Args:
		f = &builtinSubstring2ArgsSig{base}
	case tipb.ScalarFuncSig_Substring3Args:
		f = &builtinSubstring3ArgsSig{base}
	case tipb.ScalarFuncSig_SubstringIndex:
		f = &builtinSubstringIndexSig{base}
	case tipb.ScalarFuncSig_ToBase64:
		f = &builtinToBase64Sig{base, maxAllowedPacket}
	case tipb.ScalarFuncSig_Trim1Arg:
		f = &builtinTrim1ArgSig{base}
	case tipb.ScalarFuncSig_Trim2Args:
		f = &builtinTrim2ArgsSig{base}
	case tipb.ScalarFuncSig_Trim3Args:
		f = &builtinTrim3ArgsSig{base}
	case tipb.ScalarFuncSig_UnHex:
		f = &builtinUnHexSig{base}
	case tipb.ScalarFuncSig_Upper:
		f = &builtinUpperSig{base}
	case tipb.ScalarFuncSig_UpperUTF8:
		f = &builtinUpperUTF8Sig{base}
	case tipb.ScalarFuncSig_ToBinary:
		f = &builtinInternalToBinarySig{base}
	case tipb.ScalarFuncSig_FromBinary:
		// TODO: set the `cannotConvertStringAsWarning` accordingly
		f = &builtinInternalFromBinarySig{base, false}
	case tipb.ScalarFuncSig_CastVectorFloat32AsString:
		f = &builtinCastVectorFloat32AsStringSig{base}
	case tipb.ScalarFuncSig_CastVectorFloat32AsVectorFloat32:
		f = &builtinCastVectorFloat32AsVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_LTVectorFloat32:
		f = &builtinLTVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_LEVectorFloat32:
		f = &builtinLEVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_GTVectorFloat32:
		f = &builtinGTVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_GEVectorFloat32:
		f = &builtinGEVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_NEVectorFloat32:
		f = &builtinNEVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_EQVectorFloat32:
		f = &builtinEQVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_NullEQVectorFloat32:
		f = &builtinNullEQVectorFloat32Sig{base}
	case tipb.ScalarFuncSig_VectorFloat32AnyValue:
		f = &builtinVectorFloat32AnyValueSig{base}
	case tipb.ScalarFuncSig_VectorFloat32IsNull:
		f = &builtinVectorFloat32IsNullSig{base}
	case tipb.ScalarFuncSig_VecAsTextSig:
		f = &builtinVecAsTextSig{base}
	case tipb.ScalarFuncSig_VecDimsSig:
		f = &builtinVecDimsSig{base}
	case tipb.ScalarFuncSig_VecL1DistanceSig:
		f = &builtinVecL1DistanceSig{base}
	case tipb.ScalarFuncSig_VecL2DistanceSig:
		f = &builtinVecL2DistanceSig{base}
	case tipb.ScalarFuncSig_VecNegativeInnerProductSig:
		f = &builtinVecNegativeInnerProductSig{base}
	case tipb.ScalarFuncSig_VecCosineDistanceSig:
		f = &builtinVecCosineDistanceSig{base}
	case tipb.ScalarFuncSig_VecL2NormSig:
		f = &builtinVecL2NormSig{base}
	case tipb.ScalarFuncSig_FTSMatchWord:
		f = &builtinFtsMatchWordSig{base}
	default:
		e = ErrFunctionNotExists.GenWithStackByArgs("FUNCTION", sigCode)
		return nil, e
	}
	f.setPbCode(sigCode)
	return f, nil
}

func newDistSQLFunctionBySig(ctx BuildContext, sigCode tipb.ScalarFuncSig, tp *tipb.FieldType, args []Expression) (Expression, error) {
	f, err := getSignatureByPB(ctx, sigCode, tp, args)
	if err != nil {
		return nil, err
	}
	// derive collation information for string function, and we must do it
	// before doing implicit cast.
	funcName := tipb.ScalarFuncSig_name[int32(sigCode)]
	argTps := make([]types.EvalType, 0, len(args))
	for _, arg := range args {
		argTps = append(argTps, arg.GetType(ctx.GetEvalCtx()).EvalType())
	}
	ec, err := deriveCollation(ctx, funcName, args, f.getRetTp().EvalType(), argTps...)
	if err != nil {
		return nil, err
	}
	funcF := &ScalarFunction{
		FuncName: ast.NewCIStr(fmt.Sprintf("sig_%T", f)),
		Function: f,
		RetType:  f.getRetTp(),
	}
	funcF.SetCoercibility(ec.Coer)
	return funcF, nil
}

// PBToExprs converts pb structures to expressions.
func PBToExprs(ctx BuildContext, pbExprs []*tipb.Expr, fieldTps []*types.FieldType) ([]Expression, error) {
	exprs := make([]Expression, 0, len(pbExprs))
	for _, expr := range pbExprs {
		e, err := PBToExpr(ctx, expr, fieldTps)
		if err != nil {
			return nil, errors.Trace(err)
		}
		if e == nil {
			return nil, errors.Errorf("pb to expression failed, pb expression is %v", expr)
		}
		exprs = append(exprs, e)
	}
	return exprs, nil
}

// PBToExpr converts pb structure to expression.
func PBToExpr(ctx BuildContext, expr *tipb.Expr, tps []*types.FieldType) (Expression, error) {
	evalCtx := ctx.GetEvalCtx()
	switch expr.Tp {
	case tipb.ExprType_ColumnRef:
		_, offset, err := codec.DecodeInt(expr.Val)
		if err != nil {
			return nil, err
		}
		return &Column{Index: int(offset), RetType: tps[offset]}, nil
	case tipb.ExprType_Null:
		return &Constant{Value: types.Datum{}, RetType: types.NewFieldType(mysql.TypeNull)}, nil
	case tipb.ExprType_Int64:
		return convertInt(expr.Val, expr.FieldType)
	case tipb.ExprType_Uint64:
		return convertUint(expr.Val, expr.FieldType)
	case tipb.ExprType_String:
		return convertString(expr.Val, expr.FieldType)
	case tipb.ExprType_Bytes:
		return &Constant{Value: types.NewBytesDatum(expr.Val), RetType: types.NewFieldType(mysql.TypeString)}, nil
	case tipb.ExprType_MysqlBit:
		return &Constant{Value: types.NewMysqlBitDatum(expr.Val), RetType: types.NewFieldType(mysql.TypeString)}, nil
	case tipb.ExprType_Float32:
		return convertFloat(expr.Val, true)
	case tipb.ExprType_Float64:
		return convertFloat(expr.Val, false)
	case tipb.ExprType_MysqlDecimal:
		return convertDecimal(expr.Val, expr.FieldType)
	case tipb.ExprType_MysqlDuration:
		return convertDuration(expr.Val)
	case tipb.ExprType_MysqlTime:
		return convertTime(expr.Val, expr.FieldType, evalCtx.Location())
	case tipb.ExprType_MysqlJson:
		return convertJSON(expr.Val)
	case tipb.ExprType_MysqlEnum:
		return convertEnum(expr.Val, expr.FieldType)
	case tipb.ExprType_TiDBVectorFloat32:
		return convertVectorFloat32(expr.Val)
	}
	if expr.Tp != tipb.ExprType_ScalarFunc {
		panic("should be a tipb.ExprType_ScalarFunc")
	}
	// Then it must be a scalar function.
	args := make([]Expression, 0, len(expr.Children))
	for _, child := range expr.Children {
		if child.Tp == tipb.ExprType_ValueList {
			results, err := decodeValueList(child.Val)
			if err != nil {
				return nil, err
			}
			if len(results) == 0 {
				return &Constant{Value: types.NewDatum(false), RetType: types.NewFieldType(mysql.TypeLonglong)}, nil
			}
			args = append(args, results...)
			continue
		}
		arg, err := PBToExpr(ctx, child, tps)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	sf, err := newDistSQLFunctionBySig(ctx, expr.Sig, expr.FieldType, args)
	if err != nil {
		return nil, err
	}

	return sf, nil
}

func convertTime(data []byte, ftPB *tipb.FieldType, tz *time.Location) (*Constant, error) {
	ft := PbTypeToFieldType(ftPB)
	_, v, err := codec.DecodeUint(data)
	if err != nil {
		return nil, err
	}
	var t types.Time
	t.SetType(ft.GetType())
	t.SetFsp(ft.GetDecimal())
	err = t.FromPackedUint(v)
	if err != nil {
		return nil, err
	}
	if ft.GetType() == mysql.TypeTimestamp && tz != time.UTC {
		err = t.ConvertTimeZone(time.UTC, tz)
		if err != nil {
			return nil, err
		}
	}
	return &Constant{Value: types.NewTimeDatum(t), RetType: ft}, nil
}

func decodeValueList(data []byte) ([]Expression, error) {
	if len(data) == 0 {
		return nil, nil
	}
	list, err := codec.Decode(data, 1)
	if err != nil {
		return nil, err
	}
	result := make([]Expression, 0, len(list))
	for _, value := range list {
		result = append(result, &Constant{Value: value})
	}
	return result, nil
}

func convertInt(val []byte, tp *tipb.FieldType) (*Constant, error) {
	var d types.Datum
	_, i, err := codec.DecodeInt(val)
	if err != nil {
		return nil, errors.Errorf("invalid int % x", val)
	}
	d.SetInt64(i)
	return &Constant{Value: d, RetType: PbTypeToFieldType(tp)}, nil
}

func convertUint(val []byte, tp *tipb.FieldType) (*Constant, error) {
	var d types.Datum
	_, u, err := codec.DecodeUint(val)
	if err != nil {
		return nil, errors.Errorf("invalid uint % x", val)
	}
	d.SetUint64(u)
	return &Constant{Value: d, RetType: PbTypeToFieldType(tp)}, nil
}

func convertString(val []byte, tp *tipb.FieldType) (*Constant, error) {
	var d types.Datum
	d.SetBytesAsString(val, collate.ProtoToCollation(tp.Collate), uint32(tp.Flen))
	return &Constant{Value: d, RetType: PbTypeToFieldType(tp)}, nil
}

func convertFloat(val []byte, f32 bool) (*Constant, error) {
	var d types.Datum
	_, f, err := codec.DecodeFloat(val)
	if err != nil {
		return nil, errors.Errorf("invalid float % x", val)
	}
	if f32 {
		d.SetFloat32(float32(f))
	} else {
		d.SetFloat64(f)
	}
	return &Constant{Value: d, RetType: types.NewFieldType(mysql.TypeDouble)}, nil
}

func convertDecimal(val []byte, ftPB *tipb.FieldType) (*Constant, error) {
	ft := PbTypeToFieldType(ftPB)
	_, dec, precision, frac, err := codec.DecodeDecimal(val)
	var d types.Datum
	d.SetMysqlDecimal(dec)
	d.SetLength(precision)
	d.SetFrac(frac)
	if err != nil {
		return nil, errors.Errorf("invalid decimal % x", val)
	}
	return &Constant{Value: d, RetType: ft}, nil
}

func convertDuration(val []byte) (*Constant, error) {
	var d types.Datum
	_, i, err := codec.DecodeInt(val)
	if err != nil {
		return nil, errors.Errorf("invalid duration %d", i)
	}
	d.SetMysqlDuration(types.Duration{Duration: time.Duration(i), Fsp: types.MaxFsp})
	return &Constant{Value: d, RetType: types.NewFieldType(mysql.TypeDuration)}, nil
}

func convertJSON(val []byte) (*Constant, error) {
	var d types.Datum
	_, d, err := codec.DecodeOne(val)
	if err != nil {
		return nil, errors.Errorf("invalid json % x", val)
	}
	if d.Kind() != types.KindMysqlJSON {
		return nil, errors.Errorf("invalid Datum.Kind() %d", d.Kind())
	}
	return &Constant{Value: d, RetType: types.NewFieldType(mysql.TypeJSON)}, nil
}

func convertVectorFloat32(val []byte) (*Constant, error) {
	v, _, err := types.ZeroCopyDeserializeVectorFloat32(val)
	if err != nil {
		return nil, errors.Errorf("invalid VectorFloat32 %x", val)
	}
	var d types.Datum
	d.SetVectorFloat32(v)
	return &Constant{Value: d, RetType: types.NewFieldType(mysql.TypeTiDBVectorFloat32)}, nil
}

func convertEnum(val []byte, tp *tipb.FieldType) (*Constant, error) {
	_, uVal, err := codec.DecodeUint(val)
	if err != nil {
		return nil, errors.Errorf("invalid enum % x", val)
	}
	// If uVal is 0, it should return Enum{}
	var e = types.Enum{}
	if uVal != 0 {
		e, err = types.ParseEnumValue(tp.Elems, uVal)
		if err != nil {
			return nil, err
		}
	}
	d := types.NewMysqlEnumDatum(e)
	return &Constant{Value: d, RetType: FieldTypeFromPB(tp)}, nil
}
