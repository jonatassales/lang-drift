import { compileShader } from './compile-shader'
import { fragmentSource, vertexSource } from './shaders'

export function createProgram(gl: WebGL2RenderingContext): WebGLProgram | null {
	const vertex = compileShader(gl, gl.VERTEX_SHADER, vertexSource)
	const fragment = compileShader(gl, gl.FRAGMENT_SHADER, fragmentSource)
	const program = gl.createProgram()

	if (!vertex || !fragment || !program) {
		gl.deleteShader(vertex)
		gl.deleteShader(fragment)
		gl.deleteProgram(program)

		return null
	}

	gl.attachShader(program, vertex)
	gl.attachShader(program, fragment)
	gl.linkProgram(program)
	gl.deleteShader(vertex)
	gl.deleteShader(fragment)

	if (!gl.getProgramParameter(program, gl.LINK_STATUS)) {
		gl.deleteProgram(program)

		return null
	}

	return program
}
