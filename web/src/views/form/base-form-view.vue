<template>
  <a-card :style="{ width: '100%'}" title="公告列表">
    <template #extra>
      <a-button size="mini" type="primary" @click="add">
        <template #icon>
          <icon-plus/>
        </template>
        <!-- Use the default slot to avoid extra spaces -->
        <template #default>发表公告</template>
      </a-button>
    </template>
    <div class="body">
      <a-card v-for="item  of dataList" :key="item.identity" :size="mini"
              :style="{ width: '30%' ,height:'30%' , margin:'5px 20px' }" bordered
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
            <a-button size="mini" status="success" type="primary" @click="compile(item)">编辑</a-button>
            <a-button size="mini" type="primary" @click="del(item)">删除</a-button>
          </a-space>
        </template>

      </a-card>
    </div>
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
                  :action=baseURL()+upload()
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
  <ModalDialog ref="modalDialogRefs" :title="actionTitle" @confirm="onDataFormConfirms">
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
          <a-textarea v-model.trim="notice.text" placeholder="请输入公告内容" show-word-limit>
          </a-textarea>
        </a-form-item>
        <a-form-item label="公告状态">
          <a-switch v-model="off" checked-color="#165DFF" unchecked-color="#14C9C9">
            <template #checked-icon>
              <icon-check/>
            </template>
            <template #unchecked-icon>
              <icon-close/>
            </template>
          </a-switch>
          <!--          </a-textarea>-->
        </a-form-item>
        <a-form-item label="图片">
          <a-card :style="{ width: '100%',  height: '30%' }">
            <a-space>
              <a-image
                  :src="notice.url?notice.url:notice.url"
                  width="auto"
              />
              <a-upload
                  :fileList="file ? [file] : []"
                  :headers="{ 'Authorization': 'Bearer ' + token}"
                  :show-file-list="false"
                  :action=baseURL()+upload()
                  show-cancel-button
                  @change="onChange"
                  @progress="onProgress"
                  @success="success"
              >
                <template #upload-button>
                  <div
                      :class="`arco-upload-list-item${   file && file.status === 'error' ? ' arco-upload-list-item-error' : ''   }`">
                    <div v-if="file && file.url">
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
import {delNotice, getNoticeList, publishNotice, updateNotice, upload} from "@/api/url";
import {Message, Modal} from "@arco-design/web-vue";
import usePost from "@/hooks/usePost";
import useUserStore from "@/store/modules/user";
import useGet from "@/hooks/useGet";
import {baseURL} from "@/api/axios.config";

export default defineComponent({
  methods: {
    upload() {
      return upload
    },
    baseURL() {
      return baseURL
    }
  },
  components: {AddButton},
  name: 'BaseFromView',
  setup: function () {
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const modalDialogRefs = ref<ModalDialogType | null>(null)
    const actionTitle = ref('发布公告')
    const file = ref();
    const token = localStorage.getItem("token")
    const post = usePost()
    const get = useGet()
    const userStore = useUserStore()
    const dataList = ref([]) as any
    let off = ref(false)
    let id: string = ""

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
      actionTitle.value = "发布公告"
      notice.title = ''
      notice.url = ''
      notice.text = ''
      modalDialogRef.value?.toggle()
    }

    function del(item: any) {
      console.log(item)
      Modal.confirm({
        title: '提示',
        content: '确定要删除此信息，删除后不可恢复？',
        onOk() {
          get({
            url: delNotice,
            headers: {
              Authorization: "Bearer " + userStore.token
            },
            data: {
              "id": item.identity,
            }
          }).then((res) => {
            doRefresh()
            Message.success("删除公告成功！")
          }).catch(error => {
            Message.error(error.message)
          })
        },
      })
    }

    // 编辑
    function compile(item: any) {
      console.log(item.identity)
      actionTitle.value = "修改公告"
      modalDialogRefs.value?.toggle()
      notice.title = item.title
      notice.url = item.url
      notice.text = item.text
      off.value = item.status === 1
      id = item.identity
    }

    function onDataFormConfirm() {
      post({
        url: publishNotice,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: {
          "title": notice.title,
          "text": notice.text,
          "url": notice.url,
        }
      }).then((data) => {
        doRefresh()
        Message.success("发布成功")
      }).catch(error => {
        Message.error(error.message)
      })
      modalDialogRef.value?.toggle()
      notice.text = ''
      notice.url = ''
      notice.title = ''
      file.value = ''
    }

    function onDataFormConfirms() {
      post({
        url: updateNotice,
        headers: {
          Authorization: "Bearer " + userStore.token
        },

        data: {
          "title": notice.title,
          "text": notice.text,
          "url": notice.url,
          "status": off.value ? 1 : 2,
          "id": id,
        }
      }).then((res) => {
        doRefresh()
        Message.success("更新公告成功")
      }).catch(error => {
        Message.error(error.message)
      })
      modalDialogRefs.value?.toggle()
      notice.text = ''
      notice.url = ''
      notice.title = ''
    }

    function success(response: any) {
      console.log(response)
      const {url} = response.response.data[0]
      console.log(url, file.value)
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
      modalDialogRef,
      modalDialogRefs,
      actionTitle,
      notice,
      file,
      token,
      dataList,
      off,
      id,
      onDataFormConfirm,
      onDataFormConfirms,
      add,
      del,
      compile,
      success,
      onChange,
      onProgress,
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