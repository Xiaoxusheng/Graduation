<template>
  <div class="main-container">
    <div class="box-wrapper">
      <a-card
          :body-style="{ padding: '10px' }"
          :bordered="true"
          class="card-border-radius personal-box"
          hoverable
      >
        <template #actions>
          <a-space>
          <span class="icon-hover"
                style="font-size: 20px"> </span>
            <span class="icon-hover" style="font-size: 10px"> {{ "发布" }} </span>
            <a-divider direction="vertical"/>
            <a-button size="mini" status="success" type="primary" @click="compile">修改密码</a-button>
            <a-divider direction="vertical"/>
          </a-space>
        </template>
        <div class="info-wrapper">
          <div class="avatar-wrapper">
            <div
                :class="{ 'avatar-touch': touched, 'avatar-end': uploaded }"
                class="avatar"
                @mouseenter="avatarTouchStart"
            >
              <img :src="avatar" alt=""/>
            </div>
            <div class="flex items-center justify-center camera-layer" @click="uploadAvatar">
              <icon-camera style="color: #fff; font-size: 30px"/>
            </div>
          </div>
          <div class="text-xl">
          </div>
          <div class="des-wrapper">
            <i class="el-icon-edit"></i>
            冰冻三尺，非一日之寒，成大事者不拘小节。
          </div>
          <div class="text-wrapper">
            <div class="label"> 昵称：</div>
            <div class="value"> {{ info.name }}</div>
          </div>
          <div class="text-wrapper">
            <div class="label"> 工号：</div>
            <div class="value"> {{ info.uid }}</div>
          </div>
          <div class="text-wrapper">
            <div class="label"> 性别：</div>
            <div class="value"> {{ info.sex == 1 ? '女' : '男' }}</div>
          </div>
          <div class="text-wrapper">
            <div class="label"> 生日：</div>
            <div class="value"> {{ info.birthday }}</div>
          </div>
          <div class="text-wrapper">
            <div class="label"> 部门：</div>
            <div class="value"> {{ department[info.department_id] }}</div>
          </div>
          <div class="text-wrapper">
            <div class="label"> 入职时间：</div>
            <div class="value"> {{ info.create }}</div>
          </div>
          <div class="text-wrapper">
            <div class="label"> IP地址：</div>
            <div class="value"> {{ info.IP }}</div>
          </div>
          <div class="mt-4">
            <a-space :style="{ flexWrap: 'wrap' }" align="center">
              <a-tag color="green" size="small">技术控</a-tag>
              <a-tag color="green" size="small">爱学习</a-tag>
              <a-tag color="green" size="small">大嘴巴</a-tag>
              <a-tag color="green" size="small">宅男</a-tag>
            </a-space>
          </div>
          <a-card-meta title="">

          </a-card-meta>
        </div>

      </a-card>

    </div>
  </div>
  <ModalDialog ref="modalDialogRefs" :title="actionTitle" @confirm="onDataFormConfirms">
    <template #content>
      <a-form ref="formRef" :labelCol="{ span: 4 }" :model="pwd">
        <a-form-item :rules="[
            { required: true, message: '请输入旧密码' },
            { min: 3, max: 50, message: '长度在 3 - 10个字符' },
          ]" :validate-trigger="['change', 'input']" field="pwd" label="旧密码">
          <a-input v-model.trim="pwd.pwd" placeholder="请输入旧密码">
            <template #suffix>
              <icon-info-circle/>
            </template>
          </a-input>
        </a-form-item>
        <a-form-item :rules="[
            { required: true, message: '请输入新密码' },
            { min: 3, max: 50, message: '长度在 3 - 10个字符' },
          ]" :validate-trigger="['change', 'input']" field="password" label="新密码">
          <a-input v-model.trim="pwd.password" placeholder="请输入新密码">
            <template #suffix>
              <icon-info-circle/>
            </template>
          </a-input>
        </a-form-item>

      </a-form>
    </template>
  </ModalDialog>

</template>

<script lang="ts">
import useUserStore from '@/store/modules/user'
import {defineComponent, onMounted, reactive, ref} from 'vue'
import {get, Response} from "@/api/http";
import {changePassword, userinfo} from "@/api/url";
import {Message} from "@arco-design/web-vue";
import {ModalDialogType} from "@/types/components";

export default defineComponent({
  name: 'info',
  setup() {
    const touched = ref(false)
    const uploaded = ref(false)
    const avatarTouchStart = () => {
      touched.value = true
    }
    const uploadAvatar = () => {
      uploaded.value = true
      setTimeout(() => {
        touched.value = false
        uploaded.value = false
      }, 1000)
    }
    let info = reactive({
      name: '',
      birthday: '',
      sex: 1,
      uid: undefined,
      status: 1,
      department_id: 1,
      phone: '',
      position: 1,
      IP: "127.0.0.1",
      create: 0,
    })
    const department = {
      1: "程序部",
      2: "人事部",
      3: "财务部",
      4: "销售部",
      5: "法务部"
    }
    const modalDialogRefs = ref<ModalDialogType | null>(null)
    const actionTitle = ref('修改密码')
    const userStore = useUserStore()
    const pwd = reactive({
      pwd: "",
      password: ""
    })

    function getInfo() {
      // 再次请求信息接口
      get({
        url: userinfo,
        headers: {
          Authorization: "Bearer " + "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6IjBhZDM0Zjk5LWM1NGQtNDgwOC04NTI3LTUzMTJlZTlhZmIzOCIsImV4cCI6MTcxMDMwNTY0NSwiaWF0IjoxNzA5ODczNjQ1fQ.7rNYbDdy6tCwQoFyDZ7X0rUoNSlrmER-VMV1QIsf2tw"
        },
      }).then(({data}: Response) => {
        info.birthday = Getdate(data.birthday)
        info.status = data.status
        info.sex = data.sex
        info.uid = data.uid
        info.phone = data.phone
        info.department_id = data.department_id
        info.position = data.position
        info.name = data.name
        info.IP = data.IP
        info.create = GetTime(data.CreatedAt)
        console.log(info)
      }).catch(error => {
        Message.error(error.message)
      })
    }

    function Getdate(item: any) {
      const date = new Date(item);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      return `${year}-${month}-${day} `
    }

    function GetTime(item: any) {
      const date = new Date(item);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const hours = String(date.getHours()).padStart(2, "0");
      const minutes = String(date.getMinutes()).padStart(2, "0");
      const seconds = String(date.getSeconds()).padStart(2, "0");
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
    }

    function del(item: any) {
      console.log(item)
    }

    // 编辑
    function compile(item: any) {
      modalDialogRefs.value?.toggle()

      console.log(1)
    }

    function onDataFormConfirms() {
      get({
        url: changePassword,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: pwd
      }).then(({data}: Response) => {
        Message.success('修改成功')
      }).catch(error => {
        Message.error(error.message)
      })
      pwd.pwd = ""
      pwd.password = ""
      modalDialogRefs.value?.toggle()
    }

    onMounted(getInfo)
    return {
      modalDialogRefs,
      actionTitle,
      touched,
      uploaded,
      avatar: userStore.avatar,
      nickName: userStore.nickName,
      info,
      department,
      pwd,
      del,
      compile,
      Getdate,
      GetTime,
      onDataFormConfirms,
      avatarTouchStart,
      uploadAvatar,
    }
  },
})
</script>
<style lang="less" scoped>
.box-wrapper {
  .personal-box {

    .info-wrapper {
      text-align: center;

      .avatar-wrapper {
        display: inline-block;
        width: 6rem;
        height: 6rem;
        margin-top: 20px;
        position: relative;

        .avatar {
          position: absolute;
          top: 0;
          left: 0;
          right: 0;
          bottom: 0;
          transform-origin: bottom;
          transform: rotate(0deg);
          z-index: 1;
          transition: all 0.5s ease-in-out;

          & > img {
            width: 100%;
            height: 100%;
            border-radius: 50%;
            border: 2px solid rgb(245, 241, 7);
          }
        }

        .avatar-touch {
          transform: rotate(180deg);
        }

        .avatar-end {
          transform: rotate(0deg);
        }

        .camera-layer {
          position: absolute;
          top: 0;
          bottom: 0;
          left: 0;
          right: 0;
          background: rgba(0, 0, 0, 0.6);
          border-radius: 50%;
        }
      }

      .des-wrapper {
        width: 70%;
        margin: 0 auto;
        font-size: 14px;
        padding: 15px;
      }

      .text-wrapper {
        font-size: 0.8rem;
        margin-top: 20px;
        width: 50%;
        margin: 0 auto;

        .label {
          display: inline-block;
          width: 40%;
          text-align: right;
        }

        .value {
          display: inline-block;
          width: 60%;
          text-align: left;
        }
      }

      .text-wrapper + .text-wrapper {
        margin-top: 15px;
      }
    }
  }

  .message-wrapper {
    border-bottom: 1px solid #f5f5f5;
    padding-bottom: 10px;

    .notify {
      width: 10px;
      height: 10px;
      border-radius: 50%;
    }

    .message-title {
      font-size: 14px;
    }

    .content {
      font-size: 12px;
      margin-top: 10px;
      line-height: 1rem;
    }
  }

  .message-wrapper + .message-wrapper {
    margin-top: 10px;
  }

  .wating-box {
    width: 30%;
    margin-left: 10px;

    .wating-item {
      padding: 10px;
      border-bottom: 1px solid #f5f5f5;
    }
  }
}
</style>
