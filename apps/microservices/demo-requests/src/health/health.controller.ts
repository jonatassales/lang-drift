import { Controller, Get } from '@nestjs/common'

interface HealthReply {
	status: 'ok'
}

// The readiness signal the Dockerfile's HEALTHCHECK and the microruntime probe.
@Controller('health')
export class HealthController {
	@Get()
	check(): HealthReply {
		return {
			status: 'ok',
		}
	}
}
