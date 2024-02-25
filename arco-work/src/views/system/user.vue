<template>
  <a-watermark content="员工管理系统">
    <div style="width: 100%; height:100%;"/>
    <TableBody>
      <template #header>
        <TableHeader ref="tableHeaderRef" :show-filter="false">
          <template #table-config>
            <add-button @add="onAddItem"/>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-table
            :bordered="true"
            :data="dataList"
            :hoverable="true"
            :loading="tableLoading"
            :pagination="false"
            :row-selection="{ selectedRowKeys }"
            :rowKey="rowKey"
            :scroll="{ y: tableHeight }"
            :stripe="true"
            size="small"
            table-layout-fixed
            @selection-change="onSelectionChange"
        >
          <template #columns>
            <a-table-column
                v-for="item of tableColumns"
                :key="item.identity"
                :align="item.align"
                :data-index="(item.key as string)"
                :fixed="item.fixed"
                :title="(item.title as string)"
            >
              <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                {{ rowIndex + 1 }}
              </template>
              <template v-else-if="item.key === 'sex'" #cell="{ record }">
                <a-tag :color="record.sex === 1 ? 'green' : 'red'">
                  {{ record.sex === 1 ? '男' : '女' }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'position'" #cell="{ record }">
                <a-tag :color="record.position === 0 ? 'green' : 'black'">
                  {{ map[record.position] }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'department_id'" #cell="{ record }">
                <a-tag :color="record.department_id === 0 ? 'green' : 'green'">
                  {{ department[record.department_id] }}
                </a-tag>
              </template>

              <template v-else-if="item.key === 'image_Url'" #cell="{}">
                <a-avatar :size="30" :style="{ backgroundColor: 'var(--color-primary-light-1)' }">
                  <IconUser/>
                </a-avatar>
              </template>
              <template v-else-if="item.key === 'actions'" #cell="{ record }">
                <a-space>
                  <a-button size="mini" status="success" @click="onUpdateItem(record)">
                    编辑
                  </a-button>
                  <a-button size="mini" status="danger" @click="onDeleteItem(record)">
                    删除
                  </a-button>
                </a-space>
              </template>
              <template v-else-if="item.key === 'status'" #cell="{ record }">
                <a-tag v-if="record.status === 1" color="blue" size="small">正常</a-tag>
                <a-tag v-else color="red" size="small">离职</a-tag>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </template>


      <template #footer>
        <TableFooter ref="tableFooterRef" :pagination="pagination"/>
      </template>
    </TableBody>
    <!--    继续-->
    <ModalDialog ref="modalDialogRef" :title="actionTitle" @confirm="onDataFormConfirm">
      <template #content>
        {{ employerInfo }}

        <a-form ref="formRef" :labelCol="{ span: 4 }" :model="employerInfo">
          <a-form-item :rules="[
            { required: true, message: '请输入员工姓名' },
            { min: 3, max: 10, message: '长度在 3 - 10个字符' },
          ]" :validate-trigger="['change', 'input']" field="name" label="员工姓名">
            <a-input v-model="employerInfo.name" placeholder="请输入员工姓名">
              <template #suffix>
                <icon-info-circle/>
              </template>
            </a-input>
          </a-form-item>
          <a-form-item :rules="[
            { required: true, message: '请输入部门编号' },
            { min: 0, max: 2, message: '长度在 1个字符' },
          ]" :validate-trigger="['change', 'input']" field="depCode" label="部门编号">
            <a-input v-model="employerInfo.department" placeholder="请输入部门编号">
              <template #suffix>
                <icon-info-circle/>
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="电话" name="phone">
            <a-input v-model="employerInfo.phone" placeholder="请输入员工电话">
              <template #suffix>
                <icon-info-circle/>
              </template>
            </a-input>
          </a-form-item>
          <a-form-item label="员工生日" name="birthday">
            <a-date-picker
                v-model="employerInfo.birthday"
                default-value="2000-01-01"
                format=YYYY:MM:DD
                style="width: 100%"
                type="time"
            />
          </a-form-item>
          <a-form-item label="职位" name="position">
            <a-select v-model.number="employerInfo.position" placeholder="请选择职位">
              <a-option v-for="(value, key) of map" :key="key" :value="key as number">
                {{ value }}
              </a-option>
            </a-select>

          </a-form-item>
          <a-form-item label="性别" name="sex">
            <a-radio-group v-model="employerInfo.sex">
              <a-radio :value="1">男</a-radio>
              <a-radio :value="2">女</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="状态" name="status">
            <a-radio-group v-model="employerInfo.status">
              <a-radio :value="1">正常</a-radio>
              <a-radio :value="2">离职</a-radio>
            </a-radio-group>
          </a-form-item>
        </a-form>
      </template>
    </ModalDialog>
  </a-watermark>
</template>

<script lang="ts">
import {add_employer, delete_employer, employerList, update_employer} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn, useTableHeight,} from '@/hooks/table'
import {Form, Message, Modal} from '@arco-design/web-vue'
import {defineComponent, getCurrentInstance, onMounted, reactive, ref} from 'vue'
import AddButton from "@/components/AddButton.vue";
import useUserStore from "@/store/modules/user";
import {ModalDialogType} from "@/types/components";
import usePost from '@/hooks/usePost'
import useGet from "@/hooks/useGet";

export default defineComponent({
  components: {AddButton},
  name: 'UserList',
  setup: function () {
    const table = useTable()
    const rowKey = useRowKey('id')
    const pagination = usePagination(doRefresh)
    const {selectedRowKeys, onSelectionChange} = useRowSelection()
    const map = {
      1: "普通员工",
      2: "副主管",
      3: "主管",
      4: "副经理",
      5: "经理",
    }

    const department = {
      1: "程序部",
      2: "人事部",
      3: "财务部",
    }

    interface EmployerInfo {
      uid: number
      name: string
      birthday: number
      sex: number
      department: number
      status: number
      position: number | undefined
      phone: number
    }

    const formRef = ref<typeof Form>()


    const formModel = ref({})
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const actionTitle = ref('添加员工')
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '工号',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '姓名',
        key: 'name',
        dataIndex: 'name',
      },
      {
        title: '性别',
        key: 'sex',
        dataIndex: 'sex',
      },
      {
        title: '头像',
        key: 'image_Url',
        dataIndex: 'image_Url ',
      },
      {
        title: '生日',
        key: 'birthday',
        dataIndex: 'birthday',
      },
      {
        title: '在职状态',
        key: 'status',
        dataIndex: 'status',
      },
      {
        title: '职位',
        key: 'position',
        dataIndex: 'position',
      },
      {
        title: '部门',
        key: 'department_id',
        dataIndex: 'department_id',
      }, {
        title: '电话',
        key: 'phone',
        dataIndex: 'phone',
      },
      {
        title: 'IP',
        key: 'IP',
        dataIndex: 'IP',
      },
      {
        title: '操作',
        key: 'actions',
        dataIndex: 'actions',
      },
    ])
    const expandAllFlag = ref(true)
    const userStore = useUserStore()
    const post = usePost()
    const get = useGet()
    let add = true
    //
    let employerInfo = reactive<EmployerInfo>({
      uid: 0,
      name: '',
      birthday: 0,
      sex: 1,
      department: 1,
      status: 1,
      position: undefined,
      phone: 1,
    })

    function doRefresh() {
      get({
        url: employerList,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            offset: pagination.page,
            limit: pagination.pageSize,
          }
        },
      }).then((res) => {
        res.data.forEach(i => {
          const date = new Date(i.birthday);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          return i.birthday = `${year}-${month}-${day} `
        })
        table.handleSuccess(res)
        pagination.setTotalSize(table.dataList.length)
      }).catch(error => {
        console.log(error)
      })
    }


    function onDeleteItems() {
      if (selectedRowKeys.value.length === 0) {
        Message.error('请选择要删除的数据')
        return
      }
      sl


      Modal.confirm({
        title: '提示',
        content: '确定要删除此数据吗？',
        cancelText: '取消',
        okText: '删除',
        onOk: () => {
          Message.success(
              '数据模拟删除成功，所选择的Keys为：' + JSON.stringify(selectedRowKeys.value)
          )
        },
      })
    }

    function onDeleteItem(item: any) {
      Modal.confirm({
        title: '提示',
        content: '确定要删除此数据吗？',
        cancelText: '取消',
        okText: '删除',
        onOk: () => {
          console.log(item)
          get({
            url: delete_employer,
            headers: {
              Authorization: "Bearer " + userStore.token,
            },
            data: {
              uid: item.uid,
            }
          }).then((res) => {
            table.dataList.splice(table.dataList.indexOf(item), 1)

            pagination.setTotalSize(table.dataList.length)
            Message.success('删除成功')
          }).catch(error => {
            Message.error(error.toString())
            console.log(error)
          })
        },
      })
    }

    // 编辑
    function onDataFormConfirm() {
      modalDialogRef.value?.toggle()
      if (add) {
        employerInfo.birthday = Math.round(new Date(employerInfo.birthday).getTime() / 1000);
        console.log(employerInfo)
        post({
          url: add_employer,
          headers: {
            Authorization: "Bearer " + userStore.token
          },
          data: employerInfo
        }).then(() => {
          table.dataList.push(employerInfo)
          Message.success("添加成功")
        }).catch((error) => {
          Message.error(error.message)
        })
      } else {
        employerInfo.birthday = Math.round(new Date(employerInfo.birthday).getTime() / 1000);

        post({
              url: update_employer,
              headers: {
                Authorization: "Bearer " + userStore.token
              },
              data: employerInfo
            },
        ).then(() => {

          Message.success("修改成功")
        }).catch(error => {
          Message.error(error.message)
        })
      }


    }


    function onAddItem() {
      add = true
      employerInfo.phone = 1
      employerInfo.sex = 1
      employerInfo.status = 1
      employerInfo.name = ''
      employerInfo.birthday = 0
      employerInfo.position = undefined
      actionTitle.value = '添加员工'
      modalDialogRef.value?.toggle()
    }

    function onUpdateItem(record: EmployerInfo) {
      add = false
      actionTitle.value = '更新员工信息'
      modalDialogRef.value?.toggle()
      employerInfo.uid = record.uid
      employerInfo.phone = record.phone
      employerInfo.sex = record.sex
      employerInfo.status = record.status
      employerInfo.name = record.name
      employerInfo.birthday = record.birthday
      employerInfo.position = record.position
    }

    onMounted(async () => {
      table.tableHeight.value = await useTableHeight(getCurrentInstance())
      doRefresh()
    })
    return {
      ...table,
      rowKey,
      selectedRowKeys,
      onSelectionChange,
      expandAllFlag,
      tableColumns,
      pagination,
      formModel,
      actionTitle,

      modalDialogRef,
      onAddItem,
      onDeleteItem,
      department,
      employerInfo,
      formRef,
      map,
      add,
      onUpdateItem,
      onDataFormConfirm
    }
  },
})
</script>
