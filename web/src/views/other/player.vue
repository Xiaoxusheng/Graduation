<template>
  <a-card title="视频播放器（集成西瓜视频插件SDK）">
    <div class="flex justify-center">
      <div id="player" ref="playerRef"></div>
    </div>
  </a-card>
</template>

<script lang="ts">
import {defineComponent, onMounted, ref} from 'vue'
import Player from 'xgplayer'

const playUrl =
    'https://sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-720p.mp4'
const poster =
    'https://lf9-cdn-tos.bytecdntp.com/cdn/expire-1-M/byted-player-videos/1.0.0/poster.jpg'
export default defineComponent({
  name: 'Player',
  setup() {
    const playerRef = ref<HTMLDivElement>()
    onMounted(() => {
      new Player({
        el: playerRef.value,
        url: playUrl,
        poster,
        playbackRate: [0.5, 0.75, 1, 1.5, 2],
        defaultPlaybackRate: 1,
      })
    })
    return {
      playerRef,
    }
  },
})
</script>
