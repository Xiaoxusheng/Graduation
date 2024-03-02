import {defineStore} from 'pinia'
import {UserState} from '../types'
import store from '../pinia'

import Avatar from '@/assets/img_avatar.gif'
import {nextTick} from 'vue'

const defaultAvatar = Avatar

const useUserStore = defineStore('user-info', {
    state: () => {
        return {
            identity: 0,
            // roleId: 0,
            token: '',
            userName: '',
            nickName: '',
            avatar: defaultAvatar,
        }
    },
    actions: {
        saveUser(userInfo: UserState) {
            return new Promise<UserState>((resolve) => {
                this.identity = userInfo.identity
                // this.roleId = userInfo.roleId
                this.token = userInfo.token
                this.userName = userInfo.userName ? userInfo.userName : 'admin'
                this.nickName = userInfo.nickName ? userInfo.nickName : '超级管理员'
                this.avatar = userInfo.avatar || defaultAvatar
                resolve(userInfo)
                this.setToken()
            })
        },
        setToken() {
            localStorage.setItem("token", this.token)
        },
        isTokenExpire() {
            return !this.token
        },
        changeNickName(newNickName: string) {
            this.nickName = newNickName
        },
        logout() {
            return new Promise<void>((resolve) => {
                this.$reset()
                localStorage.clear()
                sessionStorage.clear()
                resolve()
            })
        },
    },
    presist: {
        enable: true,
        resetToState: true,
        option: {
            exclude: ['userName'],
        },
    },
})

export default useUserStore

/**
 * 返回一个Promise对象，目的是防止在非vue组件使用中的时候出现插件还没有出初始化完成导致持久化操作失败的bug。
 * 使用方式：
 *  import { useUserStoreContext } from '@/store/modules/user'
 useUserStoreContext().then((store) => {
 console.log(store.nickName)
 })
 * @returns Promise<typeof useUserStore>
 */
export async function useUserStoreContext() {
    await nextTick()
    return useUserStore(store)
}
