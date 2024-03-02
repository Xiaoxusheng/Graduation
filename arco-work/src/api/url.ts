import {baseURL} from './axios.config'

export const baseAddress = baseURL

export const test = '/test'

// 登录
export const login = '/login'
// 退出登录
export const logouts = '/logout'
// 获取个人信息
export const info = '/info'


export const updateUserInfo = '/updateUser'

export const addUserInfo = '/addUser'

// 菜单列表
export const getMenuListByRoleId = '/menu_list'

export const getAllMenuByRoleId = '/getAllMenuByRoleId'

export const deleteUserById = '/deleteUserById'

// 部门列表
export const getDepartmentList = '/get_department_list'
// 获取部门信息
export const getDepartmentInfoList = '/department_info'
// 添加部门
export const addDepartment = '/add_department'
// 更新部门信息
export const updateDepartment = '/update_department'
// 删除部门
export const delDepartment = '/del_department'

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
// 获取请假申请列表
export const getLeaveApplicationList = '/get_LeaveApplication_list'
//获取加班申请列表
export const getOvertimeApplicationList = '/get_overtimeApplication_list'
// 补卡申请列表
export const makeCardApplicationList = '/make_card_application_list'

// 请假审批
export const leaveApplication = '/leave_application'

// 加班审批
export const overtimeApplication = '/overtime_application'

// 补卡申请审批
export const makeCardApplication = '/make_card_application'


// 考勤
export const getAllClockIn = '/get_all_clockIn'

// 工资列表
export const getSalary = '/get_salary'


// 日志
export const logList = '/log_list'


export const getTableList = '/employer_list'

export const getCardList = '/getCardList'

export const getCommentList = '/getCommentList'

declare module 'vue' {
    interface ComponentCustomProperties {
        $urlPath: Record<string, string>
    }
}
