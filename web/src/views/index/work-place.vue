<template>
  <div class="main-container">
    <a-card
        :bodyStyle="{ padding: '10px' }"
        :bordered="false"
        :headStyle="{ padding: '0 10px' }"
        class="card-border-radius"
        size="small"
        title="工作台"
    >
      <a-row class="margin-top" wrap>
        <a-col :lg="16" :md="16" :sm="16" :xl="14" :xs="24">
          <div class="flex justify-center items-center">
            <div class="avatar-wrapper">
              <img :src="avatar"/>
            </div>
            <div class="flex flex-col justify-around ml-3.5 flex-1">
              <div class="text-lg">早上好，Andy，青春只有一次，别让自己过得不精彩</div>
              <div class="text-sm text-gray-500 mt-1">今日有小雨，出门别忘记带伞哦~</div>
            </div>
          </div>
        </a-col>
        <a-col :lg="8" :md="8" :sm="8" :xl="10" :xs="24">
          <div class="flex justify-end items-center h-full w-full mt-4">
            <div class="flex flex-col justify-around align-end item-action">
              <div class="text-gray">项目数</div>
              <div class="text-lg mt-2">12</div>
            </div>
            <div class="flex flex-col justify-around align-end item-action">
              <div class="text-gray">待办项</div>
              <div class="text-lg mt-2">3/20</div>
            </div>
            <div class="flex flex-col justify-around align-end item-action">
              <div class="text-gray">当前日期</div>
              <div class="text-lg mt-2">{{ currentDate }}</div>
            </div>
          </div>
        </a-col>
      </a-row>
    </a-card>
    <div class="mt-3"></div>
    <a-row :gutter="[20, 10]">
      <a-col
          v-for="(item, index) of fastActions"
          :key="index"
          :md="8"
          :sm="8"
          :xl="4"
          :xs="12"
          :xxl="4"
      >
        <a-card
            :bordered="false"
            class="flex flex-col items-center justify-center fast-item-wrapper"
            @click="fastActionClick(item)"
        >
          <a-space align="center" direction="vertical">
            <component :is="item.icon" :style="{ color: item.color, fontSize: '28px' }"/>
            <span class="mt-8 text-md">{{ item.title }}</span>
          </a-space>
        </a-card>
      </a-col>
    </a-row>
    <div class="mt-3"></div>
    <a-card
        :bodyStyle="{ padding: '10px' }"
        :bordered="true"
        size="small"
    >
      <a-calendar v-model="value"/>
      select: {{ value }}
    </a-card>
  </div>
</template>
<script lang="ts">
import WeiXin from '@/assets/qrcode.jpg'
import WeiXinCustom from '@/assets/custom_weixin.jpg'
import {computed, defineComponent, reactive, ref} from 'vue'
import {useRouter} from 'vue-router'
import {random} from 'lodash-es'
import useUserStore from '@/store/modules/user'

const COLORS = ['#67C23A', '#E6A23C', '#F56C6C', '#409EFF']
const date = new Date()
export default defineComponent({
  name: 'WorkPlace',
  setup() {
    const userStore = useUserStore()
    const avatar = computed(() => userStore.avatar)
    const tempWaitingItems = reactive([] as Array<any>)
    const router = useRouter()
    const fastActionClick = ({path = '/'}) => {
      router.push(path)
    }
    const value = ref(new Date());

    return {
      tempWaitingItems,
      WeiXin,
      WeiXinCustom,
      avatar,
      currentDate: date.getFullYear() + '/' + (date.getMonth() + 1) + '/' + date.getDate(),
      fastActions: [
        {
          title: '首页',
          icon: 'icon-dashboard',
          path: '/index/home',
          color: COLORS[random(0, COLORS.length)],
        },
        {
          title: '系统管理',
          path: '/systems/department',
          icon: 'icon-settings',
          color: COLORS[random(0, COLORS.length)],
        },
        {
          title: '考勤管理',
          path: '/list/table-custom',
          icon: 'icon-list',
          color: COLORS[random(0, COLORS.length)],
        },
        {
          title: '系统日志',
          path: '/log/system',
          icon: 'icon-edit',
          color: COLORS[random(0, COLORS.length)],
        },
        {
          title: '财务',
          path: '/salary/salary',
          icon: 'icon-share-alt',
          color: COLORS[random(0, COLORS.length)],
        },
        {
          title: '更多功能',
          path: '/other/chart/icons',
          icon: 'icon-apps',
          color: COLORS[random(0, COLORS.length)],
        },
      ],
      value,
      fastActionClick,
      onPreview: function (item: any) {
        window.open(item.target)
      },
    }
  },
})
</script>

<style lang="less" scoped>
.avatar-wrapper {
  width: 3rem;
  height: 3rem;
  max-width: 3rem;
  max-height: 3rem;
  min-width: 3rem;
  min-height: 3rem;

  & > img {
    width: 100%;
    height: 100%;
    border-radius: 50%;
    border: 2px solid yellowgreen;
  }
}

.item-action {
  position: relative;
  padding: 0 30px;
}

.item-action::after {
  position: absolute;
  top: 20%;
  right: 0;
  height: 60%;
  content: '';
  display: block;
  width: 1px;
  background-color: var(--border-color);
}

div.item-action:last-child::after {
  width: 0;
}

.fast-item-wrapper {
  height: 80px;
  border-radius: 8px;

  .anticon {
    font-size: 20px;
  }
}

.fast-item-wrapper:hover {
  cursor: pointer;
  box-shadow: 0px 0px 10px #ddd;
}
</style>
