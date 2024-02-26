import {AxiosResponse} from 'axios'
import {Message} from '@arco-design/web-vue'
import useUserStore from '@/store/modules/user'


export default function (response: AxiosResponse): AxiosResponse {
    const useStore = useUserStore()

    if (response.data.code === 20011) {
        Message.error('当前用户登录已过期，请重新登录')
        setTimeout(() => {
            ;(useStore as any).logout() && (useStore as any).onLogout()
        }, 1500)
    }
    return response
}
