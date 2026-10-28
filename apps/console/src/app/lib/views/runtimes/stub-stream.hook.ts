'use client'

import { useEffect, useState } from 'react'

import type { RuntimeWidgetSnapshot } from '@data/types'

import { parseStubFrame } from './stub-frame'

/** One data-only peer per visible product list; teardown owns channel and timers. */
export function useStubStream(enabled: boolean, productId?: string) {
	const [widgets, setWidgets] = useState<RuntimeWidgetSnapshot[] | null>(null)
	const [connected, setConnected] = useState(false)
	useEffect(() => {
		setWidgets(null)
		setConnected(false)
		if (!enabled) return
		const controller = new AbortController()
		let peer: RTCPeerConnection | undefined
		let channel: RTCDataChannel
		try {
			if (typeof RTCPeerConnection !== 'function') return
			peer = new RTCPeerConnection({
				iceServers: [],
			})
			channel = peer.createDataChannel('extub')
		} catch {
			peer?.close()

			return
		}
		const connection = peer
		let lastMessage = 0
		let sequence = 0
		channel.onmessage = (event) => {
			const frame = parseStubFrame(event.data, sequence)
			if (!frame) {
				setWidgets(null)
				setConnected(false)

				return
			}
			sequence = frame.sequence
			lastMessage = Date.now()
			setWidgets(
				frame.widgets.filter(
					(item) => !productId || item.runtime.productId === productId,
				),
			)
			setConnected(true)
		}
		const disconnected = () => {
			setConnected(false)
			setWidgets(null)
		}
		channel.onclose = disconnected
		channel.onerror = disconnected
		const timer = window.setInterval(() => {
			if (lastMessage && Date.now() - lastMessage > 5000) disconnected()
		}, 1000)
		async function connect() {
			try {
				const offer = await connection.createOffer()

				await connection.setLocalDescription(offer)

				await new Promise<void>((resolve, reject) => {
					const timeout = window.setTimeout(
						() => reject(new Error('ice_timeout')),
						5000,
					)
					const finish = () => {
						window.clearTimeout(timeout)
						resolve()
					}
					connection.onicegatheringstatechange = () => {
						if (connection.iceGatheringState === 'complete') finish()
					}
					controller.signal.addEventListener(
						'abort',
						() => {
							window.clearTimeout(timeout)
							reject(new Error('cancelled'))
						},
						{
							once: true,
						},
					)
					if (connection.iceGatheringState === 'complete') finish()
				})

				const response = await fetch('/api/extub/webrtc/offer', {
					method: 'POST',
					headers: {
						'content-type': 'application/json',
					},
					body: JSON.stringify(connection.localDescription),
					signal: controller.signal,
				})

				if (!response.ok) throw new Error('negotiation_failed')
				const answer = await response.json()

				await connection.setRemoteDescription(answer)
			} catch {
				if (!controller.signal.aborted) disconnected()
				connection.close()
			}
		}
		void connect()

		return () => {
			controller.abort()
			window.clearInterval(timer)
			channel.onmessage = null
			channel.onclose = null
			channel.onerror = null
			channel.close()
			connection.close()
		}
	}, [enabled, productId])

	return {
		widgets,
		connected,
	}
}
