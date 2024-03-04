<template>
  <div class="main-container">
    <a-card :content-style="{ padding: '0px' }" :header-style="{ padding: '5px' }" title="打印图片">
      <template #extra>
        <a-space>
          <a-button size="small" type="primary" @click="printImage">打印</a-button>
        </a-space>
      </template>
      <div class="image-wrapper">
        <a-card>
          <a-space :style="{ width: '100%' }" direction="vertical">
            <a-upload
                :fileList="file ? [file] : []"
                :show-file-list="false"
                action="/"
                show-cancel-button
                @change="onChange"
                @progress="onProgress"
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
      </div>
    </a-card>
    <a-card
        :content-style="{ padding: '0px' }"
        :header-style="{ padding: '5px' }"
        class="mt-4"
        title="打印HTML"
    >
      <template #extra>
        <a-button size="small" type="primary" @click="printHtml">打印</a-button>
      </template>
      <div id="htmlWrapper">
        <a-table :columns="columns" :data="dataList" :pagination="false"></a-table>
      </div>
    </a-card>
  </div>
</template>

<script lang="ts">
import printJS from 'print-js'
import imagePath from '@/assets/logo.png'
import {defineComponent, ref} from 'vue'
import {IconEdit, IconPlus} from '@arco-design/web-vue/es/icon';

export default defineComponent({
  components: {IconPlus, IconEdit},
  name: 'Print',
  setup: function () {
    const src = ref(" ")
    const file = ref();

    const onChange = (_, currentFile) => {
      file.value = {
        ...currentFile,
        // url: URL.createObjectURL(currentFile.file),
      };
    };
    const onProgress = (currentFile) => {
      file.value = currentFile;
    };

    function printImage() {
      printJS({
        documentTitle: file.value.name,
        printable: file.value.url,
        type: 'image',
        showModal: false,
        css: '',
      })
    }

    function printHtml() {
      printJS({
        printable: 'htmlWrapper',
        type: 'html' || 'pdf' || 'xsl' || 'docx',
        targetStyles: ['*'],
      })
    }

    return {
      printImage,
      printHtml,
      file,
      onChange,
      onProgress,
      src,
      imagePath,
      dataList: [
        {
          name: '张三',
          age: 10,
          gender: '男',
        },
        {
          name: '李四',
          age: 40,
          gender: '男',
        },
        {
          name: '王五',
          age: 30,
          gender: '女',
        },
      ],
      columns: [
        {
          title: '姓名',
          dataIndex: 'name',
        },
        {
          title: '年龄',
          dataIndex: 'age',
        },
        {
          title: '姓别',
          dataIndex: 'gender',
        },
      ],
    }
  },
})
</script>

<style lang="less" scoped>
.image-wrapper {
  width: 30%;
  margin: 0 auto;

  & > img {
    width: 100%;
  }
}

.html-wrapper {
  width: 80%;
  margin: 0 auto;

  & > h1 {
    color: red;
  }
}
</style>
