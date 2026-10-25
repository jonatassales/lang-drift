export const vertexSource = `#version 300 es
void main() {
  vec2 position = vec2((gl_VertexID << 1) & 2, gl_VertexID & 2);
  gl_Position = vec4(position * 2.0 - 1.0, 0.0, 1.0);
}`

// Layered, translucent petal sheets in the brand's purples; every detail is
// smooth geometry, not grain. The colors arrive as uniforms derived from the
// token roles, so one program serves both color schemes: \`scheme\` is 0 in
// dark, where light is added over a near-black field, and 1 in light, where
// lilac veils are laid over the pale one at a lower intensity.
export const fragmentSource = `#version 300 es
precision highp float;
uniform vec2 resolution;
uniform float elapsed;
uniform float scheme;
uniform vec3 background;
uniform vec3 surface;
uniform vec3 primary;
uniform vec3 primaryHover;
uniform vec3 shadeDeep;
uniform vec3 shadeSoft;
uniform vec3 glow;
uniform vec3 violet;
out vec4 outputColor;

float softBand(float distance, float width) {
  float scaled = distance / width;
  return exp(-scaled * scaled);
}

vec2 petalCoordinates(vec2 point, vec2 center, vec2 radius, float angle) {
  mat2 rotation = mat2(cos(angle), -sin(angle), sin(angle), cos(angle));
  return rotation * (point - center) / radius;
}

vec3 curledPetal(vec3 underneath, vec2 point, vec2 center, vec2 radius,
                 float angle, float phase, vec3 lightColor, vec3 ink, float opacity) {
  vec2 local = petalCoordinates(point, center, radius, angle);
  local.x += 0.24 * local.y * local.y + 0.08 * sin(local.y * 2.8 + phase);
  local.y += 0.12 * sin(local.x * 2.2 - phase);
  float edge = length(local * vec2(1.0, 1.0 + 0.09 * local.x));
  float coverage = 1.0 - smoothstep(0.90, 1.13, edge);

  // The face bends through the light; its shading varies inside the silhouette.
  float curl = local.y + 0.25 * local.x * local.x;
  vec3 normal = normalize(vec3(
    0.65 * local.x + 0.22 * sin(local.y * 3.0 + phase),
    0.90 * curl,
    0.72 + 0.30 * cos(local.x * 2.1)
  ));
  float lighting = max(0.0, dot(normal, normalize(vec3(-0.70, -0.45, 1.0))));
  float turnedFace = softBand(curl - 0.35, 0.46);
  float pigment = 0.15 + 0.40 * (1.0 - lighting) + 0.20 * turnedFace;
  vec3 face = mix(lightColor, ink, pigment);

  // A curled sheet has two faces meeting along a bent lip. Broad opposing
  // light and shadow bands give that turn depth independently of its outline.
  float fold = local.y + 0.38 * local.x + 0.32 * local.x * local.x;
  fold += 0.08 * sin(local.x * 3.0 + phase);
  float rolledFace = softBand(fold - 0.12, 0.27);
  float litFace = softBand(fold + 0.28, 0.24);
  face = mix(face, mix(lightColor, ink, 0.65), rolledFace * 0.62);
  face = mix(face, lightColor, litFace * 0.68);

  // Transmitted light belongs to the upper/left lip, not an outline of the
  // oval. It carries the glow shade, less of it in the light scheme.
  float backlitEdge = softBand(edge - 0.91, 0.21);
  backlitEdge *= 1.0 - smoothstep(-0.30, 0.65, local.x + local.y);
  vec3 rim = mix(lightColor, glow, 0.55 - 0.25 * scheme);
  face = mix(face, rim, backlitEdge * 0.64);
  float innerCrease = softBand(curl - 0.62, 0.19);
  face = mix(face, mix(lightColor, ink, 0.56), innerCrease * 0.20);

  float contact = softBand(edge - 1.09, 0.30);
  contact *= smoothstep(-0.45, 0.65, local.y);
  underneath = mix(underneath, mix(lightColor, ink, 0.49), contact * opacity * 0.28);
  return mix(underneath, face, coverage * opacity);
}

void main() {
  vec2 point = gl_FragCoord.xy / resolution;
  point.y = 1.0 - point.y;
  float aspect = resolution.x / resolution.y;
  point.x = (point.x - 0.5) * clamp(aspect / 2.23, 0.90, 1.3) + 0.5;
  float phase = elapsed * 1.04719755;
  float opening = 0.5 + 0.5 * cos(phase);
  point += vec2(
    0.025 * sin(point.y * 4.5 + phase),
    0.018 * sin(point.x * 4.0 - phase)
  );

  // The light every sheet turns toward: a lifted glow in dark, a lilac veil
  // in light. The field itself leans on the deep shade in light.
  float intensity = mix(1.0, 0.55, scheme);
  vec3 lift = mix(glow, shadeSoft, scheme);
  vec3 pale = mix(mix(background, surface, 0.24), shadeSoft, 0.18 + 0.22 * scheme);
  vec3 fieldInk = mix(primary, shadeDeep, 0.30 * scheme);

  float leftBend = point.x + 0.08 * sin(point.y * 3.8 - phase);
  float leftShade = 1.0 - smoothstep(-0.15, 0.57, leftBend);
  vec3 color = mix(pale, fieldInk, 0.10 + leftShade * 0.33);

  // Far petals remain soft behind the foreground, with independent bending
  // phases; the second one carries the violet neighbor of the brand.
  color = curledPetal(color, point, vec2(0.88, 0.95), vec2(0.43, 0.52),
                     -0.38, phase + 1.4, pale, primaryHover, 0.96);
  color = curledPetal(color, point, vec2(0.97, -0.24), vec2(0.50, 0.37),
                     0.18, phase + 2.2, pale, violet, 0.58);

  // A broad diagonal sheet enters from the lower-left. Its bow catches light
  // across a surface, while a soft trough in the deep shade supplies overlap.
  float sheetDistance = point.x + point.y * 0.58 - 0.71;
  sheetDistance += 0.075 * sin(point.y * 3.7 - phase);
  float sheetCoverage = 1.0 - smoothstep(-0.035, 0.11, sheetDistance);
  float sheetBend = softBand(sheetDistance + 0.21, 0.13);
  float sheetLight = softBand(sheetDistance + 0.055, 0.11);
  vec3 sheetColor = mix(pale, primary, 0.12 + sheetBend * 0.37);
  sheetColor = mix(sheetColor, pale, sheetLight * (0.46 + opening * 0.30));
  float trailingFold = softBand(sheetDistance + 0.40, 0.10);
  sheetColor = mix(sheetColor, mix(pale, violet, 0.40), trailingFold * 0.24);
  float trough = softBand(sheetDistance + 0.30, 0.07);
  sheetColor = mix(sheetColor, shadeDeep, trough * 0.18);
  color = mix(color, sheetColor, sheetCoverage * smoothstep(0.02, 0.40, point.y));

  // The luminous opening changes width as the sheets turn, rather than sliding
  // a rigid stripe over the completed composition.
  float transmittedLight = softBand(sheetDistance - 0.075, 0.13 + opening * 0.05);
  color = mix(color, lift, transmittedLight * (0.38 + opening * 0.40) * intensity);

  vec2 frontCenter = vec2(0.94 + 0.025 * sin(phase), 0.13 + 0.025 * cos(phase));
  vec2 frontRadius = vec2(0.28 + 0.020 * cos(phase), 0.51 + 0.020 * sin(phase));
  color = curledPetal(color, point, frontCenter, frontRadius,
                     0.50 + 0.07 * sin(phase), phase, lift,
                     mix(primaryHover, shadeDeep, 0.35), 0.96);

  // Diffuse light at the crown and base joins the separate layers without
  // grain. In dark the crown is lit by the glow instead of sinking to black.
  vec2 crown = (point - vec2(0.52, -0.05)) / vec2(0.27, 0.50);
  vec3 crownLight = mix(background, glow, 0.25 * (1.0 - scheme));
  color = mix(color, crownLight, exp(-dot(crown, crown)) * 0.48);
  vec2 lowerLight = (point - vec2(0.68, 0.86)) / vec2(0.22, 0.34);
  color = mix(color, pale, exp(-dot(lowerLight, lowerLight)) * (0.10 + opening * 0.08));

  // The canvas meets the page background along its bottom edge.
  color = mix(color, background, smoothstep(0.70, 1.0, point.y) * 0.85);
  outputColor = vec4(color, 1.0);
}`
