<template>
  <div class="space-wrap">
    <canvas ref="canvasRef" class="space-canvas"></canvas>
    <!-- Nebula blobs (CSS) -->
    <div class="nebula nebula-1"></div>
    <div class="nebula nebula-2"></div>
    <div class="nebula nebula-3"></div>
    <!-- Shooting stars -->
    <div class="shooting shooting-1"></div>
    <div class="shooting shooting-2"></div>
    <div class="shooting shooting-3"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const canvasRef = ref<HTMLCanvasElement | null>(null);
let animId = 0;

onMounted(() => {
  const canvas = canvasRef.value!;
  const ctx = canvas.getContext('2d')!;

  function resize() {
    canvas.width  = canvas.offsetWidth;
    canvas.height = canvas.offsetHeight;
  }
  resize();

  const ro = new ResizeObserver(resize);
  ro.observe(canvas);

  // Stars
  interface Star { x: number; y: number; r: number; o: number; s: number; d: 1 | -1 }
  const stars: Star[] = Array.from({ length: 220 }, () => ({
    x: Math.random() * canvas.width,
    y: Math.random() * canvas.height,
    r: Math.random() * 1.3 + 0.15,
    o: Math.random(),
    s: Math.random() * 0.012 + 0.003,
    d: Math.random() > 0.5 ? 1 : -1,
  }));

  function draw() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    stars.forEach((s) => {
      s.o += s.s * s.d;
      if (s.o >= 1)    { s.o = 1;    s.d = -1; }
      if (s.o <= 0.04) { s.o = 0.04; s.d = 1;  }

      // Bright stars: soft glow
      if (s.r > 1) {
        const grd = ctx.createRadialGradient(s.x, s.y, 0, s.x, s.y, s.r * 3);
        grd.addColorStop(0, `rgba(220, 240, 255, ${s.o * 0.6})`);
        grd.addColorStop(1, `rgba(220, 240, 255, 0)`);
        ctx.beginPath();
        ctx.arc(s.x, s.y, s.r * 3, 0, Math.PI * 2);
        ctx.fillStyle = grd;
        ctx.fill();
      }

      ctx.beginPath();
      ctx.arc(s.x, s.y, s.r, 0, Math.PI * 2);
      ctx.fillStyle = `rgba(220, 235, 255, ${s.o})`;
      ctx.fill();
    });

    animId = requestAnimationFrame(draw);
  }
  draw();

  onUnmounted(() => {
    cancelAnimationFrame(animId);
    ro.disconnect();
  });
});
</script>

<style scoped>
.space-wrap {
  position: absolute;
  inset: 0;
  overflow: hidden;
  background: #03081a;
}

.space-canvas {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
}

/* ─── Nebula blobs ───────────────────────── */
.nebula {
  position: absolute;
  border-radius: 50%;
  filter: blur(70px);
  pointer-events: none;
}
.nebula-1 {
  width: 500px;
  height: 300px;
  top: -60px;
  left: 10%;
  background: radial-gradient(ellipse, rgba(100, 60, 180, 0.45) 0%, transparent 70%);
  animation: nebula-drift 18s ease-in-out infinite alternate;
}
.nebula-2 {
  width: 600px;
  height: 400px;
  top: 20px;
  right: -80px;
  background: radial-gradient(ellipse, rgba(20, 80, 160, 0.4) 0%, transparent 70%);
  animation: nebula-drift 24s ease-in-out infinite alternate-reverse;
}
.nebula-3 {
  width: 400px;
  height: 250px;
  bottom: -40px;
  left: 30%;
  background: radial-gradient(ellipse, rgba(60, 130, 200, 0.3) 0%, transparent 70%);
  animation: nebula-drift 20s ease-in-out infinite alternate;
}
@keyframes nebula-drift {
  0%   { transform: translate(0, 0) scale(1);       opacity: 0.7; }
  50%  { transform: translate(20px, -15px) scale(1.08); opacity: 1; }
  100% { transform: translate(-15px, 10px) scale(0.95); opacity: 0.5; }
}

/* ─── Shooting stars ─────────────────────── */
.shooting {
  position: absolute;
  width: 80px;
  height: 1px;
  background: linear-gradient(90deg, rgba(255,255,255,0) 0%, rgba(255,255,255,0.9) 60%, rgba(255,255,255,0) 100%);
  transform: rotate(-35deg);
  opacity: 0;
}
.shooting-1 {
  top: 15%;
  left: -80px;
  animation: shoot 8s ease-in-out 1s infinite;
}
.shooting-2 {
  top: 35%;
  left: -80px;
  width: 120px;
  animation: shoot 12s ease-in-out 5s infinite;
}
.shooting-3 {
  top: 50%;
  left: -80px;
  width: 60px;
  animation: shoot 9s ease-in-out 9s infinite;
}
@keyframes shoot {
  0%   { transform: rotate(-35deg) translateX(0);    opacity: 0; }
  5%   { opacity: 1; }
  25%  { transform: rotate(-35deg) translateX(110vw); opacity: 0; }
  100% { transform: rotate(-35deg) translateX(110vw); opacity: 0; }
}
</style>
