<template>
  <a-watermark content="员工管理系统">
    <div style="width: 100%; height:80%;"/>
    <TableBody>
      <template #header>
        <TableHeader ref="tableHeaderRef" :show-filter="false">
          <template #search-content :title="'操作日志'">
            <a-form :model="{}" layout="inline">
              <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
                <template v-if="item.render">
                  <FormRender :formItem="item" :render="item.render"/>
                </template>
                <template v-else>
                  <template v-if="item.type === 'input'">
                    <a-input v-model="item.value.value" :placeholder="item.placeholder"/>
                  </template>
                  <template v-if="item.type === 'select'">
                    <a-select
                        v-model="item.value.value"
                        :placeholder="item.placeholder"
                        style="width: 150px"
                    >
                      <a-option
                          v-for="optionItem of item.optionItems"
                          :key="optionItem.value"
                          :value="optionItem.value"
                      >
                        {{ optionItem.label }}
                      </a-option>
                    </a-select>
                  </template>
                  <template v-if="item.type === 'date'">
                    <a-date-picker v-model="item.value.value"/>
                  </template>
                  <template v-if="item.type === 'time'">
                    <a-time-picker v-model="item.value.value" value-format="HH:mm:ss"/>
                  </template>
                  <template v-if="item.type === 'check-group'">
                    <a-checkbox-group v-model="item.value.value">
                      <a-checkbox v-for="it of item.optionItems" :key="it.value" :value="it.value">
                        {{ item.label }}
                      </a-checkbox>
                    </a-checkbox-group>
                  </template>
                </template>
              </a-form-item>
            </a-form>
          </template>
        </TableHeader>
      </template>
      <template #default>
        <a-table
            :bordered="{ wrapper: true, cell: true }"
            :column-resizable="true"
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

              <template v-else-if="item.key === 'http_code'" #cell="{ record }">
                <a-tag :color="record.http_code ==200 ? 'green' : 'red'">
                  {{ "record.http_code === 200" ? '成功' : '失败' }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'path'" #cell="{ record }">
                {{ record.path }}
              </template>
              <template v-else-if="item.key === 'time'" #cell="{ record }">
                <a-tag :color="record.time <=100 ? 'blue' : 'red'">
                  {{ record.time + 'ms' }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'method'" #cell="{ record }">
                <a-tag :color="record.method ==='GET' ? 'purple' : 'brown'">
                  {{ record.method }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'image_Url'" #cell="{}">
                <a-avatar :size="30" :style="{ backgroundColor: 'var(--color-primary-light-1)' }">
                  <IconUser/>
                </a-avatar>
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
  </a-watermark>
</template>

<script lang="ts">
import {add_employer, delete_employer, logList, update_employer} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn, useTableHeight,} from '@/hooks/table'
import {Form, Input, Message, Modal} from '@arco-design/web-vue'
import {defineComponent, getCurrentInstance, h, onMounted, reactive, ref} from 'vue'
import AddButton from "@/components/AddButton.vue";
import useUserStore from "@/store/modules/user";
import {FormItem, ModalDialogType} from "@/types/components";
import usePost from '@/hooks/usePost'
import useGet from "@/hooks/useGet";
import FormRender from "@/components/FormRender";
import type {Dayjs} from "dayjs";

export default defineComponent({
  components: {FormRender, AddButton},
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

    const conditionItems: Array<FormItem> = [
      {
        key: 'name',
        label: '用户姓名',
        type: 'input',
        placeholder: '请输入用户姓名',
        value: ref(''),
        reset: function () {
          this.value.value = ''
        },
        render: (formItem: FormItem) => {
          return h(Input, {
            placeholder: '请输入姓名',
            modelValue: formItem.value.value,
            'onUpdate:modelValue': (value) => {
              formItem.value.value = value
            },
          })
        },
      },
      {
        key: 'date',
        label: '创建日期',
        type: 'date',
        value: ref<Dayjs>(),
      },
      {
        key: 'sex',
        label: '用户姓别',
        value: ref(),
        type: 'select',
        placeholder: '请选择用户姓别',
        optionItems: [
          {
            label: '男',
            value: 1,
          },
          {
            label: '女',
            value: 2,
          },
        ],
        reset: function () {
          this.value.value = undefined
        },
      },
    ]

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
        title: '操作者',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '请求方式',
        key: 'method',
        dataIndex: 'method',
      },
      {
        title: '路径',
        key: 'path',
        dataIndex: 'path',
      },
      {
        title: 'IP',
        key: 'ip',
        dataIndex: 'ip ',
      },
      {
        title: '耗时',
        key: 'time',
        dataIndex: 'time',
      },
      {
        title: '状态',
        key: 'http_code',
        dataIndex: 'http_code',
      },
      {
        title: '操作时间',
        key: 'CreatedAt',
        dataIndex: 'CreatedAt',
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
        url: logList,
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
          const date = new Date(i.CreatedAt);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          i.CreatedAt = `${year}-${month}-${day}  ${date.getHours() > 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() > 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() > 10 ? date.getSeconds() : '0' + date.getSeconds()}`
        })
        table.handleSuccess(res)
        console.log(res)
        pagination.setTotalSize(res.count)
      }).catch(error => {
        console.log(error)
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
      formRef.value
          ?.validate()
          .then((error: any) => {
            if (error) {
              return
            }
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
          })
          .catch((error: any) => {
            console.log('error', error)
            return
          })
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
      onDataFormConfirm,
      conditionItems
    }
  },
})
</script>
