<template>
  <a-card :style="{ width: '100%'}" title="公告列表">
    <template #extra>
    </template>
    <div class="body">
      <a-card v-for="item  of dataList" :key="item.identity" :style="{ width: '30%' ,height:'30%' , margin:'5px 20px' }"
              bordered
              hoverable>
        <template #cover>
          <div
              :style="{
          height: 'auto',
          overflow: 'hidden',
        }"
          >
            <img
                :src="item.url"
                :style="{ width: '100%', height:'auto', transform: 'translateY(0px)' }"
                alt="dessert"
            />
          </div>
        </template>
        <a-card-meta style="font-size: 12px" title="">
          <template #description>
            <span style="font-size: 14px;font-weight: bolder ">  {{ item.title }}：</span>
            <br>
            <br>
            &nbsp;&nbsp;&nbsp;{{ item.text }}
          </template>
        </a-card-meta>
        <template #actions>
          <a-space>
          <span class="icon-hover"
                style="font-size: 10px"> {{
              new Date(item.date * 1000).toLocaleDateString() + " " + new Date(item.date * 1000).toLocaleTimeString()
            }} </span>
            <span class="icon-hover" style="font-size: 10px"> {{ item.uid + "发布" }} </span>
            <a-divider direction="vertical"/>

            <a-link disabled href="" style="font-size: 10px">{{ item.status == 1 ? '正常' : '下架' }}</a-link>
            <a-divider direction="vertical"/>
          </a-space>
        </template>

      </a-card>
    </div>
  </a-card>


</template>

<script lang="ts">

import {defineComponent, onMounted, ref} from "vue";
import AddButton from "@/components/AddButton.vue";
import {getNotice,} from "@/api/url";
import {Message} from "@arco-design/web-vue";
import useUserStore from "@/store/modules/user";
import useGet from "@/hooks/useGet";

export default defineComponent({
  components: {AddButton},
  name: 'Notice',
  setup: function () {
    const get = useGet()
    const userStore = useUserStore()
    const dataList = ref([]) as any

    function doRefresh() {
      get({
        url: getNotice,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
      }).then((res) => {
        dataList.value.length = 0
        res.data.forEach((i: any) => {
          dataList.value.push(i)
        })
      }).catch(error => {
        Message.error(error.message)
      })
    }

    onMounted(doRefresh)
    return {
      dataList,
    }
  }
})

</script>


<style lang="less">
.body {
  margin: 0 10px;
  display: flex;
  justify-content: flex-start;
  align-items: flex-start;
  flex-wrap: wrap;
}
</style>