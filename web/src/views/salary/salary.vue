<template>
  <div style="width: 100%; height:100%;"/>
  <TableBody>
    <template #header>
      <TableHeader ref="tableHeaderRef" :show-filter="false">
        <template #table-config :title="'工资列表'">
          <a-form :model="{}" layout="inline">
            <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
              <template v-if="item.render">
                <FormRender :formItem="item" :render="item.render"/>
              </template>
              <template v-else>
                <template v-if="item.type === 'date'">
                  <a-month-picker v-model="item.value.value"/>
                </template>
              </template>
            </a-form-item>
          </a-form>

          <a-button size="small" type="outline" @click="getSalarys">查询</a-button>
        </template>
      </TableHeader>
    </template>
    <template #default>
      <a-table
          :column-resizable="false"
          :data="dataList"
          :hoverable="true"
          :loading="tableLoading"
          :pagination="false"
          :rowKey="rowKey"
          :scroll="{ y: tableHeight }"
          :stripe="true"
          size="small"
          table-layout-fixe
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
              :width="item.width"
          >
            <template v-if="item.key === 'index'" #cell="{ rowIndex }">
              {{ rowIndex + 1 }}
            </template>
            <template v-else-if="item.key === 'expected_attendance_hours'" #cell="{ record }">
              <a-tag :color="'blue'">
                {{ record.expected_attendance_hours.toFixed(1) }}
              </a-tag>
            </template>
            <template v-else-if="item.key === 'attendance_hours'" #cell="{ record }">
              <a-tag :color="'green'">
                {{ record.attendance_hours.toFixed(1) }}
              </a-tag>
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
</template>

<script lang="ts">
import {delete_employer, getSalary,} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn, useTableHeight,} from '@/hooks/table'
import {Form, Input, Message, Modal} from '@arco-design/web-vue'
import {defineComponent, getCurrentInstance, h, onMounted, ref} from 'vue'
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
    const conditionItems: Array<FormItem> = [
      {
        key: 'uid',
        label: '员工工号',
        type: 'input',
        placeholder: '请输入员工工号',
        value: ref(''),
        reset: function () {
          this.value.value = ''
        },
        render: (formItem: FormItem) => {
          return h(Input, {
            placeholder: '请输入员工工号',
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

    ]
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
        title: '月份',
        key: 'date',
        dataIndex: 'date',
      },
      {
        title: '早退迟到次数',
        key: 'count',
        dataIndex: 'count',
        ellipsis: true,
      },
      {
        title: '总时长',
        key: 'expected_attendance_hours',
        dataIndex: 'expected_attendance_hours ',
      },
      {
        title: '上班时长',
        key: 'attendance_hours',
        dataIndex: 'attendance_hours',
      },
      {
        title: '公司补贴',
        key: 'subsidy',
        dataIndex: 'subsidy',
      }, {
        title: '带薪休假时长',
        key: 'paid_leave',
        dataIndex: 'paid_leave',
      },
      {
        title: '其他',
        key: 'other',
        dataIndex: 'other',
      },
      {
        title: '生成时间',
        key: 'CreatedAt',
        dataIndex: 'CreatedAt',
        width: 200,
      },
      {
        title: '工资',
        key: 'total',
        dataIndex: 'total',
      },
    ])
    const expandAllFlag = ref(true)
    const userStore = useUserStore()
    const post = usePost()
    const get = useGet()


    function doRefresh() {
      get({
        url: getSalary,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            date: pagination.page,
            uid: pagination.pageSize,
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

    function getSalarys() {
      let uid: any
      let date: any
      conditionItems.forEach(i => {
        console.log(i.value.value)
        if (i.key == "uid") {
          uid = i.value.value
        }
        if (i.key == "date") {
          date = i.value.value
        }
      })
      get({
        url: getSalary,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
        data: () => {
          return {
            date: date,
            uid: uid
          }
        },
      }).then((res) => {
        const date = new Date(res.data.CreatedAt);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, "0");
        const day = String(date.getDate()).padStart(2, "0");
        res.data.CreatedAt = `${year}-${month}-${day}  ${date.getHours() > 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() > 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() > 10 ? date.getSeconds() : '0' + date.getSeconds()}`
        table.handle(res.data)
        pagination.setTotalSize(1)
      }).catch(error => {
        Message.error(error.toString())
        console.log(error)
      })
      conditionItems.forEach(i => {
        i.value.value = ''
      })
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
      getSalarys,
      modalDialogRef,
      onDeleteItem,
      formRef,
      conditionItems
    }
  },
})
</script>
