<template>
  <a-card :style="{ width: '100%'}" class="arco-card-body" title="公告列表">
    <template #extra>
      <a-button size="mini" type="primary" @click="add">
        <template #icon>
          <icon-plus/>
        </template>
        <!-- Use the default slot to avoid extra spaces -->
        <template #default>发表公告</template>
      </a-button>
    </template>
    <a-card v-for="item  of dataList" :key="item.identity" :style="{ width: '20%'  , margin:'0 20px'}" hoverable>
      <template #cover>
        <div
            :style="{
          height: 'auto',
          overflow: 'hidden',
        }"
        >
          <img
              :src="item.url"
              :style="{ width: '100%', transform: 'translateY(0px)' }"
              alt="dessert"
          />
        </div>
      </template>
      <a-card-meta :style="{ size:'16px'}" title="">
        <template #description>
          <span style="font-size: 16px">  {{ item.title }}：</span>
          <br>
          <br>
          &nbsp;&nbsp;{{ item.text }}
        </template>
      </a-card-meta>
      <template #actions>
        <a-space>
          <span class="icon-hover"
                style="font-size: 10px"> {{
              new Date(item.date * 1000).toLocaleDateString() + new Date(item.date * 1000).toLocaleTimeString()
            }} </span>

          <span class="icon-hover" style="font-size: 10px"> {{ item.uid + "发布" }} </span>
          <a-divider direction="vertical"/>
          <a-button size="mini" status="success" type="primary">编辑</a-button>
          <a-button size="mini" type="primary">删除</a-button>
        </a-space>
      </template>

    </a-card>
  </a-card>
  <ModalDialog ref="modalDialogRef" :title="actionTitle" @confirm="onDataFormConfirm">
    <template #content>
      <a-form ref="formRef" :labelCol="{ span: 4 }" :model="notice">
        <a-form-item :rules="[
            { required: true, message: '请输入公告标题' },
            { min: 3, max: 50, message: '长度在 3 - 50个字符' },
          ]" :validate-trigger="['change', 'input']" field="title" label="标题">
          <a-input v-model.trim="notice.title" placeholder="请输入公告标题">
            <template #suffix>
              <icon-info-circle/>
            </template>
          </a-input>
        </a-form-item>
        <a-form-item :rules="[
            { required: true, message: '请输入公告内容' },
            { min: 10, max: 200, message: '长度在 10-200个字符' },
          ]" :validate-trigger="['change', 'input']" field="text" label="公告内容">
          <a-textarea v-model.trim="notice.text" placeholder="请输入公告内容">
          </a-textarea>
        </a-form-item>
        <a-form-item label="图片">
          <a-card :style="{ width: '30%',  height: '30%' }">
            <a-space direction="vertical">
              <a-upload
                  :fileList="file ? [file] : []"
                  :headers="{ 'Authorization': 'Bearer ' + token}"
                  :show-file-list="false"
                  action="http://127.0.0.1:8084/user/upload"
                  show-cancel-button
                  @change="onChange"
                  @progress="onProgress"
                  @success="success"
              >
                <template #upload-button>
                  <div
                      :class="`arco-upload-list-item${   file && file.status === 'error' ? ' arco-upload-list-item-error' : ''   }`">
                    <div v-if="file && file.url">
                      <a-image
                          v-if="file.url!=''"
                          :src="file.url"
                          width="auto"
                      />
                      <div class="arco-upload-list-picture-mask">
                        <IconEdit/>
                      </div>
                      <a-progress
                          v-if="file.status === 'uploading' && file.percent < 100"
                          :percent="file.percent"
                          type="circle"
                      />
                    </div>
                    <div v-else class="arco-upload-picture-card">
                      <div>
                        <IconPlus/>
                        <div>Upload</div>
                      </div>
                    </div>
                  </div>
                </template>
              </a-upload>
            </a-space>
          </a-card>
        </a-form-item>

      </a-form>
    </template>
  </ModalDialog>

</template>

<script lang="ts">

import {defineComponent, onMounted, reactive, ref} from "vue";
import {ModalDialogType} from "@/types/components";
import AddButton from "@/components/AddButton.vue";
import {getNoticeList, publishNotice} from "@/api/url";
import {Message} from "@arco-design/web-vue";
import usePost from "@/hooks/usePost";
import useUserStore from "@/store/modules/user";
import useGet from "@/hooks/useGet";

export default defineComponent({
  components: {AddButton},
  name: 'BaseFromView',
  setup: function () {
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const actionTitle = ref('发布公告')
    const file = ref();
    const token = localStorage.getItem("token")
    const post = usePost()
    const get = useGet()
    const userStore = useUserStore()
    const dataList = ref([]) as any

    interface Notice {
      title: string
      url: string
      text: string
    }


    const notice = reactive<Notice>({
      title: "",
      url: "",
      text: "",
    })

    function add() {
      modalDialogRef.value?.toggle()
    }

    function onDataFormConfirm() {
      post({
        url: publishNotice,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: {
          "title": notice.title,
          "text": notice.url,
          "url": notice.url,
        }
      }).then((data) => {
        Message.success("发布成功")
      }).catch(error => {
        Message.error(error.message)
      })
      modalDialogRef.value?.toggle()
      notice.text = ''
      notice.url = ''
      notice.title = ''
    }

    function success(response: any) {
      console.log(response)
      const {url} = response.response.data[0]
      console.log(url, file.value)
      // file.value.url = url
      notice.url = url
    }

    const onChange = (_, currentFile) => {
      file.value = {
        ...currentFile,
        // url: URL.createObjectURL(currentFile.file),
      };
    };
    const onProgress = (currentFile) => {
      file.value = currentFile;
    };

    function doRefresh() {
      get({
        url: getNoticeList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
      }).then((res) => {
        res.data.forEach((i: any) => {
          dataList.value.push(i)
        })
        console.log(dataList.value)
        Message.success("获取公告成功")
      }).catch(error => {
        Message.error(error.message)
      })
    }

    onMounted(doRefresh)

    return {
      modalDialogRef,
      actionTitle,
      notice,
      file,
      token,
      dataList,
      onDataFormConfirm,
      add,
      success,
      onChange,
      onProgress,
    }
  }

})

</script>


<style lang="less">

</style>