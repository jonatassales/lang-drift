import { BaseRpcContext } from '@nestjs/microservices'

export type RedisStreamsContextArgs = [
	stream: string,
	entryId: string,
	group: string,
	deliveries: number,
]

// What a handler can learn about the entry it is handling.
export class RedisStreamsContext extends BaseRpcContext<RedisStreamsContextArgs> {
	getStream(): string {
		return this.args[0]
	}

	getEntryId(): string {
		return this.args[1]
	}

	getGroup(): string {
		return this.args[2]
	}

	getDeliveryCount(): number {
		return this.args[3]
	}
}
