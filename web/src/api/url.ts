import {baseURL} from './axios.config'

export const baseAddress = baseURL

export const test = '/test'

// 登录
export const login = '/admin/login'
// 退出登录
export const logouts = '/admin/logout'
// 获取个人信息
export const info = '/admin/info'
// 获取员工信息
export const employerInfo = '/admin/employer_info'

// 分配账号
export const assignedAccount = '/admin/assigned_account'


export const updateUserInfo = '/admin/updateUser'

export const addUserInfo = '/admin/addUser'

// 菜单列表
export const getMenuListByRoleId = '/admin/menu_list'
// 增加菜单
export const addMenu = '/admin/add_menu'
// 删除菜单
export const delMenu = '/admin/del_menu'
// 更新菜单
export const updateMenu = '/admin/update_menu'


export const getAllMenuByRoleId = '/getAllMenuByRoleId'

export const deleteUserById = '/deleteUserById'

// 部门列表
export const getDepartmentList = '/admin/get_department_list'
// 获取部门信息
export const getDepartmentInfoList = '/admin/department_info'
// 添加部门
export const addDepartment = '/admin/add_department'
// 更新部门信息
export const updateDepartment = '/admin/update_department'
// 删除部门
export const delDepartment = '/admin/del_department'

export const getRoleList = '/getRoleList'

export const getMenuList = '/getMenuList'

export const getParentMenuList = '/getParentMenuList'

// 员工列表
export const employerList = '/admin/employer_list'
// 添加员工
export const add_employer = '/admin/add_employer'
// 删除员工
export const delete_employer = '/admin/delete_employer'
// 更新员工信息
export const update_employer = '/admin/update_employer'
// 获取请假申请列表
export const getLeaveApplicationList = '/admin/get_LeaveApplication_list'
//获取加班申请列表
export const getOvertimeApplicationList = '/admin/get_overtimeApplication_list'
// 补卡申请列表
export const makeCardApplicationList = '/admin/make_card_application_list'

// 请假审批
export const leaveApplication = '/admin/leave_application'

// 加班审批
export const overtimeApplication = '/admin/overtime_application'

// 补卡申请审批
export const makeCardApplication = '/admin/make_card_application'


// 考勤
export const getAllClockIn = '/admin/get_all_clockIn'

// 工资列表
export const getSalaryList = '/admin/get_salary_list'
// 获取某个员工所有
export const getSalary = '/admin/get_salary'
//添加工资考勤
export const salaryInfo = '/admin/salary_info'



// 管理员日志
export const logList = '/admin/log_list'
// 用户日志
export const userLogList = '/admin/user_log_list'
// 发布公告
export const publishNotice = '/admin/publish_notice'
// 更新公告
export const updateNotice = '/admin/update_notice'
// 获取公告列表
export const getNoticeList = '/admin/get_notice_list'


export const getTableList = '/employer_list'

export const getCardList = '/getCardList'

export const getCommentList = '/getCommentList'
// 用户端

// 请假申请
export const leave = '/user/leave'
// 补卡申请
export const markCard_application = '/user/markCard_application'
// 加班申请
export const overtime = '/user/overtime'
//打卡
export const clockIn = '/user/clockIn'//打卡
//考勤
export const getClockInLog = '/user/get_clockIn_log'
// 获取用户信息
export const userinfo = '/user/info'
// 修改密码
export const changePassword = '/user/change_password'
//获取工资
export const getUserSalary = '/user/get_salary'

// 图片上传
export const upload = '/user/upload'

// 加班申请


export const getNotice = '/user/get_notice_list'













declare module 'vue' {
    interface ComponentCustomProperties {
        $urlPath: Record<string, string>
    }
}
