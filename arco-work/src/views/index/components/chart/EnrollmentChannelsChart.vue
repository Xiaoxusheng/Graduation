<template>
  <div class="chart-item-container">
    <div ref="channelsChart" class="chart-item"></div>
  </div>
</template>

<script lang="ts">
import useEcharts from '@/hooks/useEcharts'
import {defineComponent, nextTick, onBeforeUnmount, onMounted, ref} from 'vue'
import {dispose} from 'echarts/core'
import {get,} from "@/api/http";
import {getDepartmentInfoList} from '@/api/url'

import useUserStore from "@/store/modules/user";

const userStore = useUserStore()

export default defineComponent({
  name: 'EnrollmentChannelsChart',
  setup() {
    const loading = ref(true)
    const channelsChart = ref<HTMLDivElement | null>(null)
    const data: any = []
    get({
      url: getDepartmentInfoList,
      headers: {
        Authorization: "Bearer " + userStore.token
      },
    },).then((res) => {
      const storedMapString = localStorage.getItem('departmentMap');
      const storedMapArray = JSON.parse(storedMapString);
      const storedMap = new Map(storedMapArray);
      console.log(storedMap)
      res.data.forEach((i: any) => {
        console.log(i)
        const d = {value: i.count, name: storedMap.get(i.department_id)}
        data.push(d)
      })
    }).catch(error => {
      console.log(error)
    })
    const init = () => {
      const option = {
        legend: {
          right: '10%',
          y: 'center',
          icon: 'circle',
          orient: 'vertical',
          formatter: function (name: string) {
            // 添加
            let total = 0
            let target = 0
            for (let i = 0; i < data.length; i++) {
              total += data[i].value
              if (data[i].name === name) {
                target = data[i].value
              }
            }
            var arr = ['{a|' + name + '}', '{b|' + ((target / total) * 100).toFixed(2) + '%}']
            return arr.join('  ')
          },
          textStyle: {
            // 添加
            rich: {
              a: {
                fontSize: 12,
                color: 'var(--color-text-2)',
              },
              b: {
                fontSize: 12,
                color: 'rgb(var(--primary-1))',
                fontWeight: 'bold',
              },
            },
          },
        },
        series: [
          {
            name: '访问来源',
            type: 'pie',
            center: ['30%', '50%'],
            radius: ['50%', '70%'],
            avoidLabelOverlap: false,
            itemStyle: {
              borderColor: '#fff',
              borderWidth: 2,
            },
            emphasis: {
              label: {
                show: true,
                fontSize: '16',
                fontWeight: 'bold',
              },
            },
            label: {
              show: false,
              position: 'center',
            },
            labelLine: {
              show: false,
            },
            data,
          },
        ],
      }
      setTimeout(() => {
        loading.value = false
        nextTick(() => {
          useEcharts(channelsChart.value as HTMLDivElement).setOption(option)
        })
      }, 1000)
    }
    const updateChart = () => {
      useEcharts(channelsChart.value as HTMLDivElement).resize()
    }
    onMounted(init)
    onBeforeUnmount(() => {
      dispose(channelsChart.value as HTMLDivElement)
    })
    return {
      loading,
      channelsChart,
      updateChart,
    }
  },
})
</script>

<style lang="less" scoped>
.chart-item-container {
  width: 100%;

  .chart-item {
    height: 100%;
  }
}
</style>
