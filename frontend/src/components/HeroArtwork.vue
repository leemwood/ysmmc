<script setup lang="ts">
import { onMounted, ref } from 'vue'

type Point = [x: number, y: number]
type Position = [x: number, height: number, z: number]
type Size = [width: number, height: number, depth: number]
type Colors = [top: string, front: string, right: string]

const canvas = ref<HTMLCanvasElement>()

onMounted(() => {
  const drawingContext = canvas.value?.getContext?.('2d')
  if (!drawingContext) return
  const context = drawingContext
  context.scale(2, 2)

  const project = ([x, height, z]: Position): Point => [
    276 + x * 17 - z * 14,
    254 + x * 7 + z * 8 - height * 17,
  ]

  function polygon([first, ...rest]: [Point, ...Point[]], color: string) {
    context.beginPath()
    context.moveTo(...first)
    for (const point of rest) context.lineTo(...point)
    context.closePath()
    context.fillStyle = color
    context.fill()
  }

  function face(color: string, first: Position, ...vertices: Position[]) {
    polygon([project(first), ...vertices.map(project)], color)
  }

  function box([x, y, z]: Position, [w, h, d]: Size, [top, front, right]: Colors) {
    const r = x + w
    const t = y + h
    const f = z + d
    face(front, [x, y, f], [r, y, f], [r, t, f], [x, t, f])
    face(right, [r, y, f], [r, y, z], [r, t, z], [r, t, f])
    face(top, [x, t, z], [r, t, z], [r, t, f], [x, t, f])
  }

  function frontPatch([x, y, z]: Position, w: number, h: number, color: string) {
    face(color, [x, y, z], [x + w, y, z], [x + w, y + h, z], [x, y + h, z])
  }

  function cross(x: number, y: number, size: number, color: string) {
    const arm = Math.max(1, Math.round(size / 3))
    context.fillStyle = color
    context.fillRect(x - size, y - arm, size * 2 + 1, arm * 2 + 1)
    context.fillRect(x - arm, y - size, arm * 2 + 1, size * 2 + 1)
  }

  // Evening sky, pixel clouds and floating blocks.
  const sun = context.createRadialGradient(375, 114, 30, 375, 114, 128)
  sun.addColorStop(0, '#ffcf90')
  sun.addColorStop(1, '#f5ab73')
  context.beginPath()
  context.arc(375, 122, 106, 0, Math.PI * 2)
  context.fillStyle = sun
  context.fill()

  polygon([[0, 275], [33, 275], [33, 260], [98, 260], [98, 275], [133, 275], [133, 316], [0, 316]], '#f47d57')
  polygon([[429, 253], [457, 253], [457, 228], [507, 228], [507, 241], [534, 241], [534, 260], [560, 260], [560, 300], [429, 300]], '#d7583f')
  polygon([[377, 87], [408, 87], [408, 76], [450, 76], [450, 87], [469, 87], [469, 97], [377, 97]], '#ffe5bb')
  polygon([[55, 127], [81, 127], [81, 116], [109, 116], [109, 127], [134, 127], [134, 137], [55, 137]], '#fba77c')
  polygon([[440, 173], [471, 173], [471, 162], [495, 162], [495, 173], [526, 173], [526, 183], [440, 183]], '#ffc490')

  const rock: Colors = ['#c78b65', '#9a544c', '#703745']
  const grass: Colors = ['#d7b963', '#a27c4f', '#775449']
  const leaves: Colors = ['#f3ca72', '#d2894e', '#ad5545']
  box([-7.8, 2.9, 2.6], [1.7, 1.1, 1.7], rock)
  box([-7.8, 4, 2.6], [1.7, 0.35, 1.7], grass)
  box([6.3, 1.5, -3.6], [1.5, 1.2, 1.5], rock)
  box([6.3, 2.7, -3.6], [1.5, 0.35, 1.5], grass)
  cross(147, 90, 8, '#ffeabf')
  cross(462, 213, 6, '#fff1cf')
  cross(192, 161, 3, '#ffd79d')
  cross(440, 118, 4, '#fff0d1')
  cross(113, 246, 4, '#ffc08b')
  cross(388, 316, 3, '#ffd099')

  // Layered island, with patches on the visible faces.
  const shadow = context.createRadialGradient(280, 347, 12, 280, 347, 130)
  shadow.addColorStop(0, '#8b332c65')
  shadow.addColorStop(1, '#8b332c00')
  context.save()
  context.translate(0, 306)
  context.scale(1, 0.30)
  context.translate(0, -210)
  context.fillStyle = shadow
  context.fillRect(140, 212, 280, 270)
  context.restore()

  box([-2.35, -3.55, -1.65], [5.3, 1.6, 4.6], ['#80504c', '#814648', '#623444'])
  box([-3.65, -2.15, -2.1], [7.8, 1.4, 5.9], rock)
  box([-4.4, -0.95, -2.7], [9, 0.95, 7.1], ['#ba7c56', '#b17452', '#8d4c45'])
  box([-4.4, 0, -2.7], [9, 0.42, 7.1], grass)
  frontPatch([-4.4, -0.38, 4.4], 1.7, 0.38, grass[1])
  frontPatch([-1.8, -0.31, 4.4], 0.9, 0.31, grass[1])
  frontPatch([1.25, -0.48, 4.4], 1.2, 0.48, grass[1])
  frontPatch([-3.15, -1.86, 3.8], 1.5, 0.37, '#b36a53')
  frontPatch([0.35, -1.4, 3.8], 0.88, 0.43, '#bd795b')
  face('#754143', [4.15, -1.7, 0.5], [4.15, -1.7, -1.45], [4.15, -1.2, -1.45], [4.15, -1.2, 0.5])

  // A small tree behind the player, and the player's contact shadow.
  box([-3.35, 0.42, -1.9], [0.45, 2.65, 0.48], ['#a06b48', '#895338', '#653f3d'])
  box([-4.18, 2.2, -2.5], [2.3, 0.93, 1.82], leaves)
  box([-3.8, 3.13, -2.22], [1.48, 0.63, 1.35], ['#ffe09b', '#eaae5f', '#c67443'])
  face('#ad8b50', [-2.1, 0.425, -0.8], [2.8, 0.425, -0.8], [2.8, 0.425, 1.7], [-2.1, 0.425, 1.7])

  // Shoes, cargo trousers and a teal jacket.
  const shoes: Colors = ['#fff2d9', '#ffe5c8', '#c5a39a']
  const jeans: Colors = ['#5d638b', '#4d527b', '#343854']
  const jacket: Colors = ['#82cabc', '#508e89', '#356a73']
  const cuffs: Colors = ['#a9d6b7', '#8fbfaa', '#5c9291']
  const skin: Colors = ['#fbd3a6', '#f5bc91', '#d78e73']
  box([-1.67, 0.42, -0.82], [1.45, 0.72, 2.25], shoes)
  box([0.23, 0.42, -0.43], [1.45, 0.72, 2.28], shoes)
  box([-1.61, 1.12, -0.9], [1.35, 2.76, 1.82], jeans)
  box([0.29, 1.12, -0.52], [1.35, 2.76, 1.87], jeans)
  frontPatch([0.5, 2.58, 1.35], 0.8, 0.7, '#686787')
  box([-2.62, 3.2, -0.88], [0.9, 0.83, 1.5], skin)
  box([-2.67, 4.02, -0.92], [0.99, 0.39, 1.61], cuffs)
  box([-2.63, 4.4, -1.02], [0.96, 2.63, 1.72], jacket)
  box([-1.71, 3.75, -1.04], [3.47, 3.6, 2.14], jacket)
  frontPatch([-1.68, 3.75, 1.1], 3.42, 0.38, cuffs[1])
  frontPatch([-0.12, 4.12, 1.1], 0.18, 2.6, '#b8d3b7')
  frontPatch([-1.32, 4.55, 1.1], 0.83, 0.27, '#78aba2')
  box([1.76, 2.98, -0.22], [1.0, 0.83, 1.62], skin)
  box([1.72, 3.8, -0.26], [1.09, 0.39, 1.72], cuffs)
  box([1.7, 4.18, -0.39], [1.13, 2.75, 1.82], jacket)
  box([-1.19, 7.04, -0.81], [2.45, 0.52, 2.07], ['#dfa164', '#eaa363', '#b86153'])
  frontPatch([0.57, 6.07, 1.26], 0.61, 1.16, '#edb176')

  // Hair, eyes and headphones are drawn on the head's pixel grid.
  const hair: Colors = ['#645064', '#43374f', '#312d45']
  box([-2.0, 7.52, -1.78], [4.02, 3.76, 3.63], skin)
  box([-2.20, 10.59, -1.99], [4.47, 1.01, 3.96], hair)
  frontPatch([-2.2, 10.26, 1.97], 1.85, 0.36, hair[1])
  frontPatch([-2.2, 9.91, 1.97], 0.63, 0.36, hair[1])
  frontPatch([1.41, 10.18, 1.97], 0.85, 0.43, hair[1])
  face('#725769', [-1.62, 11.605, -1.78], [0.75, 11.605, -1.78], [0.75, 11.605, -1.21], [-1.62, 11.605, -1.21])
  face('#7b626e', [-1.55, 11.605, -0.87], [-0.5, 11.605, -0.87], [-0.5, 11.605, -0.58], [-1.55, 11.605, -0.58])
  frontPatch([-1.36, 9.05, 1.85], 0.84, 0.53, '#fff5dc')
  frontPatch([0.56, 9.05, 1.85], 0.84, 0.53, '#fff5dc')
  frontPatch([-0.99, 9.05, 1.85], 0.38, 0.53, '#413650')
  frontPatch([0.93, 9.05, 1.85], 0.38, 0.53, '#413650')
  frontPatch([-1.56, 8.68, 1.85], 0.53, 0.23, '#e5977f')
  frontPatch([1.13, 8.68, 1.85], 0.52, 0.23, '#e5977f')
  frontPatch([-0.27, 8.38, 1.85], 0.49, 0.16, '#ad675c')
  box([2.02, 8.71, -0.28], [0.37, 1.89, 1.4], ['#f3bc7c', '#f6ae65', '#b96854'])
  box([2.39, 9.03, -0.03], [0.13, 1.12, 0.91], ['#fff0b3', '#f6cb82', '#e3a06d'])

  box([3.37, 0.43, 2.52], [0.63, 0.54, 0.73], leaves)
  box([2.98, 0.43, 3.21], [0.28, 0.37, 0.28], leaves)
})
</script>

<template>
  <div class="hero-art overflow-hidden rounded-[1.125rem]" aria-hidden="true">
    <canvas ref="canvas" width="1120" height="720" class="block h-auto w-full" />
  </div>
</template>

<style scoped>
.hero-art {
  background:
    radial-gradient(circle at 26% 16%, #ffb180, transparent 57%),
    linear-gradient(140deg, #ed835b 10%, #e46a49 63%, #cf543a 100%);
}
</style>
