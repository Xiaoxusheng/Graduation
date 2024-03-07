<template>
  <a-card title="员工申请">
    <template #extra>
      <a-space>
        <a-button size="mini" status="danger" @click="resetForm"> 重置</a-button>
        <a-button size="mini" type="primary" @click="submit"> 提交</a-button>
      </a-space>
    </template>
    <div class="form-wrapper">
      <a-form :label-col-props="{ span: 3 }" :model="{}">
        <a-form-item
            v-for="item of formItems"
            :key="item.key"
            :label="item.label"
            :row-class="[item.required ? 'form-item__require' : 'form-item__no_require']"
        >
          <template v-if="item.type === 'input'">
            <a-input v-model="item.value.value" :placeholder="item.placeholder"></a-input>
          </template>
          <template v-if="item.type === 'textarea'">
            <a-textarea
                v-model="item.value.value"
                :auto-size="{ minRows: 2, maxRows: 5 }"
                :placeholder="item.placeholder"
            />
          </template>
          <template v-if="item.type === 'img'">
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
          </template>
          <template v-if="item.type === 'select'">
            <a-select v-model="item.value.value" :placeholder="item.placeholder">
              <a-option v-for="opt of item.optionItems" :key="opt.value" :value="opt.value">
                {{ opt.label }}
              </a-option>
            </a-select>
          </template>
          <template v-if="item.type === 'date-range'">
            <a-range-picker
                v-model="item.value.value"
                showTime
            />
          </template>
        </a-form-item>
      </a-form>
    </div>
  </a-card>
</template>

<script lang="ts">
import {FormItem} from '@/types/components'
import {Message} from '@arco-design/web-vue'
import {defineComponent, ref} from 'vue'
import type {Dayjs} from 'dayjs'
import {leave, markCard_application, overtime} from "@/api/url";
import useUserStore from "@/store/modules/user";
import usePost from '@/hooks/usePost'
import useGet from "@/hooks/useGet";

export default defineComponent({
  name: 'BaseFormView',
  setup() {
    const dataForm = ref()
    const formItems = [
      {
        label: '申请类型',
        key: 'leader',
        required: true,
        disabled: ref(true),
        type: 'select',
        placeholder: '请选择申请类型',
        value: ref(undefined),
        optionItems: [
          {
            label: '加班申请',
            value: 1,
          },
          {
            label: '补卡申请',
            value: 2,
          },
          {
            label: '请假申请',
            value: 4,
          },
        ],
        validator: function () {
          if (!this.value.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
        reset: function () {
          this.value.value = undefined
        },
      },
      {
        label: '起止日期',
        key: 'startEndDate',
        type: 'date-range',
        disabled: ref(true),
        value: ref<Dayjs[]>([]),
        reset: function () {
          this.value.value = []
        },
      },
      {
        label: '说明',
        key: 'content',
        disabled: ref(true),
        type: 'textarea',
        placeholder: '请输入说明,补卡，请假都需要说明，否则不予通过',
        value: ref(null),
      },
      {
        label: '图片',
        key: 'url',
        required: true,
        placeholder: '请上传图片',
        type: 'img',
        value: ref(null),
        validator: function () {
          if (!file.value) {
            Message.error(this.placeholder || '')
            return false
          }
          return true
        },
      },
    ] as FormItem[]
    const submitLoading = ref(false)
    const userStore = useUserStore()
    const file = ref();
    const post = usePost()
    const get = useGet()
    const token = typeof localStorage !== 'undefined' && localStorage.getItem('token');
    let imgurl: any

    const onChange = (_, currentFile) => {
      file.value = {
        ...currentFile,
        // url: URL.createObjectURL(currentFile.file),
      };
    };
    const onProgress = (currentFile) => {
      file.value = currentFile;
    };

    function submit() {
      if (formItems.every((it) => (it.validator ? it.validator() : true))) {
        submitLoading.value = true
        // 加班申请
        if (formItems[0].value.value == 1) {
          post({
            url: overtime,
            headers: {
              Authorization: "Bearer " + userStore.token
            },
            data: {
              "start_time": Math.floor(new Date(formItems[1].value.value[0]).getTime() / 1000) - new Date(formItems[1].value.value[0]).getTimezoneOffset() * 60,
              "end_time": Math.floor(new Date(formItems[1].value.value[1]).getTime() / 1000) - new Date(formItems[1].value.value[1]).getTimezoneOffset() * 60,
              "url": imgurl,
            }
          }).then((data) => {
            Message.success("申请成功")
          }).catch(error => {
            Message.error(error.message)
          })
          formItems.forEach(i => {
            i.value.value = ''
          })
          file.value.url = ''

        }
        // 补卡申请
        if (formItems[0].value.value == 2) {
          if (formItems[2].value.value == "") {
            Message.error("原因不能为空")
            return
          }

          const date = new Date(formItems[1].value.value[0]);
          date.setHours(0, 0, 0, 0)
          post({
            url: markCard_application,
            headers: {
              Authorization: "Bearer " + userStore.token
            },
            data: {
              "date": Math.floor((date.getTime() / 1000)),
              "reason": formItems[2].value.value,
              "url": imgurl,
            },
          }).then((data) => {
            Message.success("申请成功")
          }).catch(error => {
            Message.error(error.message)
          })
          formItems.forEach(i => {
            i.value.value = ''
          })
          file.value.url = ''

        }
        // 请假申请
        if (formItems[0].value.value == 4) {
          if (formItems[2].value.value == null) {
            Message.error("原因不能为空")
            return
          }
          post({
            url: leave,
            headers: {
              Authorization: "Bearer " + userStore.token
            },
            data: {
              "start_time": Math.floor(new Date(formItems[1].value.value[0]).getTime() / 1000) - new Date(formItems[1].value.value[0]).getTimezoneOffset() * 60,
              "end_time": Math.floor(new Date(formItems[1].value.value[1]).getTime() / 1000) - new Date(formItems[1].value.value[1]).getTimezoneOffset() * 60,
              "reason": formItems[2].value.value,
              "url": imgurl,
            }
          }).then((data) => {
            Message.success("申请成功")
          }).catch(error => {
            Message.error(error.message)
          })
          formItems.forEach(i => {
            i.value.value = ''
          })
          file.value.url = ''
        }
        // setTimeout(() => {
        //   submitLoading.value = false
        //   Message.success(
        //       '提交成功，参数为：' +
        //       JSON.stringify(
        //           formItems.reduce((pre, cur) => {
        //             ;(pre as any)[cur.key] = (cur as any).value.value
        //             return pre
        //           }, {})
        //       )
        //   )
        // }, 3000)
      }
    }

    function resetForm() {
      formItems.forEach((it) => {
        it.reset ? it.reset() : (it.value.value = '')
      })
    }


    function success(response: any) {
      const {url} = response.response.data[0]
      console.log(url)
      file.value.url = url
      imgurl = url
    }

    return {
      dataForm,
      formItems,
      submitLoading,
      file,
      token,
      imgurl,
      onChange,
      onProgress,
      submit,
      resetForm,
      success
    }
  },
})
</script>

