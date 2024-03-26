<template>
  <div class="vaw-avatar-container">
    <a-dropdown size="large" trigger="hover" @select="handleSelect">
      <div class="action-wrapper">
        <div class="avatar">
          <a-avatar :size="30">
            <img :src="userStore.avatar"/>
          </a-avatar>
        </div>
        <span class="nick-name">
          {{ userStore.role }}
        </span>
        <icon-caret-down class="tip"/>
      </div>
      <template #content>
        <a-doption v-for="item of options" :key="item.key" :value="item.key">
          <template #icon>
            <component :is="item.icon"/>
          </template>
          {{ item.label }}
        </a-doption>
      </template>
    </a-dropdown>
  </div>
</template>

<script lang="ts">
import {Message, Modal} from '@arco-design/web-vue'
import {defineComponent} from 'vue'
import {
  IconCaretDown as CaretDownOutlined,
  IconPoweroff as LogoutOutlined,
  IconUser as UserOutlined,
} from '@arco-design/web-vue/es/icon'
import useUserStore from '@/store/modules/user'
import {useRouter} from 'vue-router'
import useGet from "@/hooks/useGet";
import {logouts} from "@/api/url";

export default defineComponent({
  name: 'VAWavatar',
  components: {UserOutlined, LogoutOutlined, CaretDownOutlined},
  setup() {
    const userStore = useUserStore()
    const get = useGet()
    const options = [
      {
        label: '个人中心',
        key: 'personal-center',
        icon: 'UserOutlined',
      },
      {
        label: '退出登录',
        key: 'logout',
        icon: 'LogoutOutlined',
      },
    ]
    const router = useRouter()

    function personalCenter() {
      router.push('/personal/info')
    }

    function logout() {
      Modal.confirm({
        title: '提示',
        content: '是否要退出当前账号？',
        okText: '退出',
        cancelText: '再想想',
        onOk: () => {
          get({
            url: logouts,
            headers: {
              Authorization: "Bearer " + userStore.token
            },
          }).then((res) => {
            Message.success("退出登录成功")
          }).catch((error) => {
            Message.error(error.message)
          })
          userStore.logout().then(() => {
            window.location.reload()
          })
        },
      })
    }

    function handleSelect(key: string) {
      switch (key) {
        case 'personal-center':
          personalCenter()
          break
        case 'logout':
          logout()
          break
      }
    }

    return {
      userStore,
      options,
      handleSelect,
    }
  },
})
</script>

<style lang="less" scoped>
.vaw-avatar-container {
  .action-wrapper {
    display: flex;
    align-items: center;

    .avatar {
      display: flex;
      align-items: center;

      & > img {
        border: 1px solid #f6f6f6;
        width: 100%;
        height: 100%;
        object-fit: cover;
        border-radius: 50%;
      }
    }

    .nick-name {
      margin: 0 5px;

      .tip {
        transform: rotate(0);
        transition: transform @transitionTime;
        margin-left: 2px;
      }
    }
  }
}

.vaw-avatar-container:hover {
  cursor: pointer;
  color: var(--primary-color);

  .tip {
    transform: rotate(180deg);
    transition: transform @transitionTime;
  }
}
</style>
