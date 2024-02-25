import {baseURL} from './axios.config'

export const baseAddress = baseURL

export const test = '/test'

// 登录
export const login = '/login'

export const updateUserInfo = '/updateUser'

export const addUserInfo = '/addUser'

// 菜单列表
export const getMenuListByRoleId = '/menu_list'

export const getAllMenuByRoleId = '/getAllMenuByRoleId'

export const deleteUserById = '/deleteUserById'

// 部门列表
export const getDepartmentList = '/get_department_list'

export const addDepartment = '/addDepartment'

export const getRoleList = '/getRoleList'

export const getMenuList = '/getMenuList'

export const getParentMenuList = '/getParentMenuList'

// 员工列表
export const employerList = '/employer_list'
// 添加员工
export const add_employer = '/add_employer'
// 删除员工
export const delete_employer = '/delete_employer'
// 更新员工信息
export const update_employer = '/update_employer'


export const getTableList = '/employer_list'

export const getCardList = '/getCardList'

export const getCommentList = '/getCommentList'

declare module 'vue' {
    interface ComponentCustomProperties {
        $urlPath: Record<string, string>
    }
}
