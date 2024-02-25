import {DeviceType} from '@/store/types'

export const projectName = '员工管理系统'

export default {
    theme: 'light',
    sideTheme: 'white',
    themeColor: '#165dff',
    projectName,
    layoutMode: 'ltr',
    sideWidth: 210,
    pageAnim: 'opacity',
    isFixedNavBar: true,
    deviceType: DeviceType.PC,
    isCollapse: false,
    flexMainHeight: false,
    mainHeight: document.body.clientHeight,
    actionBar: {
        isShowSearch: true,
        isShowMessage: true,
        isShowRefresh: true,
        isShowFullScreen: true,
    },
}
