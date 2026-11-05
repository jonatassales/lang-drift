import Ajv2020, { type AnySchema, type ValidateFunction } from 'ajv/dist/2020'

// A validator built from one of the contract's own JSON Schema files, the
// schema the CLI's tests read too, so the console's tests check against the
// contract itself rather than against the console's Zod mirror of it.
export function contractValidator(schema: AnySchema): ValidateFunction {
	const ajv = new Ajv2020({
		strictTypes: false,
	})

	return ajv.compile(schema)
}
