import { Module } from '@nestjs/common'

import { ArchiveController } from './archive.controller.js'
import { ArchiveService } from './archive.service.js'
import { archiveBucketToken } from './archive-bucket.token.js'
import { createObjectStore } from './object-store.factory.js'
import { objectStoreToken } from './object-store.token.js'
import { readArchiveBucket } from './read-archive-bucket.factory.js'

@Module({
	controllers: [ArchiveController],
	providers: [
		{
			provide: objectStoreToken,
			useFactory: createObjectStore,
		},
		{
			provide: archiveBucketToken,
			useFactory: readArchiveBucket,
		},
		ArchiveService,
	],
})
export class ArchiveModule {}
