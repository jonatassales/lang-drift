import { aiAvatars } from '@pack/react/ui/ai-avatars'
import type { VisionPoint } from '@pack/react/ui/product-vision-curve'
import type { SinapsiGraphDocument, SinapsiNode } from 'sinapsi'

type DemoTranslator = (
	key: string,
	values?: Record<string, string | number>,
) => string

// The existing illustrative scenario is preserved; these values do not define
// the Product Vision formula. Locale changes only affect its presentation.
export function getHeroPoints(
	t: DemoTranslator,
	locale: string,
): VisionPoint[] {
	const month = new Intl.DateTimeFormat(locale, {
		month: 'short',
		timeZone: 'UTC',
	})
	const date = new Intl.DateTimeFormat(locale, {
		day: 'numeric',
		month: 'short',
		timeZone: 'UTC',
	})
	const at = (monthIndex: number, day = 1) =>
		new Date(Date.UTC(2026, monthIndex, day))
	const ana = {
		initials: 'AN',
		name: 'Ana',
		src: aiAvatars.ana,
		team: t('teams.product'),
	}
	const carlos = {
		initials: 'CA',
		name: 'Carlos',
		src: aiAvatars.carlos,
		team: t('teams.platform'),
	}

	return [
		{
			label: month.format(at(3)),
			value: 91,
			event: {
				actionHref: '#movement-details',
				actionLabel: t('events.baseline.action'),
				actors: [
					{
						initials: 'MR',
						name: 'Marina Reis',
						src: aiAvatars.marina,
						team: t('teams.leadership'),
					},
				],
				classification: 'baseline',
				date: date.format(at(3, 2)),
				decision: t('events.baseline.decision'),
				delta: 0,
				id: 'baseline',
				productArea: t('areas.vision'),
				reason: t('events.baseline.reason'),
				title: t('events.baseline.title'),
			},
		},
		{ label: month.format(at(4)), value: 88 },
		{
			label: month.format(at(5)),
			value: 84,
			event: {
				actionHref: '#movement-details',
				actionLabel: t('events.pricing.action'),
				actors: [ana, carlos],
				classification: 'intentional',
				date: date.format(at(5, 28)),
				decision: t('events.pricing.decision'),
				delta: -9,
				id: 'pricing',
				productArea: t('areas.pricing'),
				reason: t('events.pricing.reason'),
				title: t('events.pricing.title'),
			},
		},
		{
			label: month.format(at(6)),
			value: 79,
			event: {
				actionHref: '#movement-details',
				actionLabel: t('events.authentication.action'),
				actors: [carlos],
				classification: 'unexplained',
				date: date.format(at(6, 22)),
				decision: t('events.authentication.decision'),
				delta: -6,
				id: 'authentication',
				productArea: t('areas.authentication'),
				reason: t('events.authentication.reason'),
				title: t('events.authentication.title'),
			},
		},
		{
			label: month.format(at(7)),
			value: 73,
			event: {
				actionHref: '#movement-details',
				actionLabel: t('events.exports.action'),
				actors: [ana],
				classification: 'review',
				date: date.format(at(7, 20)),
				decision: t('events.exports.decision'),
				delta: -3,
				id: 'exports',
				productArea: t('areas.exports'),
				reason: t('events.exports.reason'),
				title: t('events.exports.title'),
			},
		},
	]
}

// Four connected clusters expose the context of the existing four events.
// A context node selects its related event; 48 nodes do not imply 48 new scores.
// IDs, event associations and edges stay unchanged across locale/theme updates.
const graphClusters = [
	{
		eventId: 'baseline',
		avatar: 'marina',
		relations: [
			[
				'vision',
				['target', 'loop', 'strategy', 'baseline', 'scope', 'pricing'],
			],
			['target', ['outcomes', 'alignment', 'segmentation']],
			['loop', ['intent', 'quality']],
			['strategy', ['roadmap', 'principles']],
			['baseline', ['target', 'intent']],
			['scope', ['roadmap', 'exports']],
			['quality', ['principles', 'authentication']],
			['intent', ['outcomes', 'decision']],
			['outcomes', ['alignment', 'reports']],
			['principles', ['alignment', 'security']],
			['alignment', ['roadmap']],
			['roadmap', ['release']],
		],
	},
	{
		eventId: 'pricing',
		avatar: 'ana',
		relations: [
			[
				'pricing',
				[
					'decision',
					'context',
					'research',
					'packaging',
					'enterprise',
					'authentication',
				],
			],
			['decision', ['approval', 'tradeoff']],
			['context', ['discovery', 'commercial']],
			['research', ['feedback', 'segmentation']],
			['feedback', ['discovery', 'experience']],
			['segmentation', ['enterprise']],
			['packaging', ['commercial', 'tradeoff', 'formats']],
			['enterprise', ['approval', 'permissions']],
			['commercial', ['approval']],
			['discovery', ['tradeoff']],
			['approval', ['review']],
			['tradeoff', ['compatibility']],
		],
	},
	{
		eventId: 'authentication',
		avatar: 'carlos',
		relations: [
			[
				'authentication',
				[
					'onboarding',
					'experience',
					'identity',
					'permissions',
					'security',
					'exports',
				],
			],
			['onboarding', ['account', 'migration']],
			['experience', ['accessibility', 'session']],
			['accessibility', ['verification']],
			['identity', ['account', 'roles']],
			['session', ['security', 'reliability']],
			['permissions', ['roles', 'verification']],
			['security', ['verification']],
			['account', ['migration', 'retention']],
			['roles', ['migration', 'workflow']],
			['migration', []],
			['verification', ['tests']],
		],
	},
	{
		eventId: 'exports',
		avatar: 'ana',
		relations: [
			[
				'exports',
				['evidence', 'review', 'delivery', 'formats', 'reports', 'vision'],
			],
			['evidence', ['tests', 'reliability']],
			['review', ['compatibility', 'retention']],
			['delivery', ['release', 'workflow']],
			['tests', ['compatibility']],
			['reliability', ['workflow']],
			['release', ['retention', 'reports']],
			['formats', ['compatibility']],
			['reports', ['workflow']],
			['retention', ['compatibility']],
			['workflow', []],
			['compatibility', []],
		],
	},
] as const

export function getHeroGraph(
	t: DemoTranslator,
	locale: string,
): SinapsiGraphDocument {
	const events = getHeroPoints(t, locale).flatMap((point) =>
		point.event ? [point.event] : [],
	)
	const formatDelta = new Intl.NumberFormat(locale, {
		signDisplay: 'exceptZero',
	})

	return {
		graph: graphClusters.flatMap((cluster) => {
			const event = events.find((item) => item.id === cluster.eventId)
			if (!event) return []

			return cluster.relations.map(
				([id, links]): SinapsiNode => ({
					id,
					name: t(`nodes.${id}`),
					payload: {
						kind: 'drift-context',
						eventId: event.id,
					},
					links: links.map((target) => ({
						id: target,
						name: t(`nodes.${target}`),
					})),
					presentation: {
						type: 'card',
						avatarUrl: `/avatars/${cluster.avatar}.jpg`,
						avatarAlt: '',
						title: event.actors.map((actor) => actor.name).join(' + '),
						description: t('graph.cardDescription', {
							context: t(`nodes.${id}`),
							movement: event.title,
						}),
						reference: `${event.date} · ${t(`classifications.${event.classification}`)}`,
						badge: t('graph.cardBadge', {
							delta: formatDelta.format(event.delta),
						}),
					},
				}),
			)
		}),
	}
}

/** Only a click on a known node with its recorded event mapping may select a point. */
export function getHeroPointFromNodeClick(
	detail: unknown,
	graph: SinapsiGraphDocument,
	points: readonly VisionPoint[],
): number | undefined {
	if (typeof detail !== 'object' || detail === null) return undefined
	if (!('event' in detail) || detail.event !== 'click') return undefined
	if (!('id' in detail) || typeof detail.id !== 'string') return undefined
	if (
		!('payload' in detail) ||
		typeof detail.payload !== 'object' ||
		detail.payload === null
	) {
		return undefined
	}
	if (
		!('eventId' in detail.payload) ||
		typeof detail.payload.eventId !== 'string'
	) {
		return undefined
	}

	const eventId = detail.payload.eventId
	const node = graph.graph.find((item) => item.id === detail.id)
	if (!node || node.payload.eventId !== eventId) return undefined
	const pointIndex = points.findIndex((point) => point.event?.id === eventId)
	return pointIndex < 0 ? undefined : pointIndex
}
