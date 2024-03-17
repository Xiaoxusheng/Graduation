<template>
  <div class="main-container">
    <TableBody ref="tableBody">
      <template #header>
        <TableHeader
            :show-filter="false"
            title="员工申请列表"
            @search="onSearch"
            @reset-search="onResetSearch"
        >
          <template #search-content>
            <a-form :model="{}" layout="inline">
              <a-form-item v-for="item of conditionItems" :key="item.key" :label="item.label">
                <template v-if="item.render">
                  <FormRender :formItem="item" :render="item.render"/>
                </template>
                <template v-else>
                  <template v-if="item.type === 'time'">
                    <a-range-picker v-model="item.value.value" show-time/>
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
            :columns="tableColumns"
            :data="dataList"
            :loading="tableLoading"
            :pagination="false"
            :row-selection="{ selectedRowKeys, showCheckedAll }"
            :rowKey="rowKey"
            @selection-change="onSelectionChange"
        >
          <template #columns>
            <a-table-column
                v-for="item of tableColumns"
                :key="item.key"
                :align="item.align"
                :data-index="(item.key as string)"
                :fixed="item.fixed"
                :title="(item.title as string)"
                :width="item.width"
            >
              <template v-if="item.key === 'index'" #cell="{ rowIndex }">
                {{ rowIndex + 1 }}
              </template>
              <template v-else-if="item.key === 'department_id'" #cell="{ record }">
                <a-tag :color="record.department_id === 1 ? 'green' : 'blue'">
                  {{ department[record.department_id] }}
                </a-tag>
              </template>
              <template v-else-if="item.key === 'avatar'" #cell="{ record }">
                <a-avatar
                    :autocapitalize="30"
                    :style="{ backgroundColor: 'var(--color-primary-light-1)' }"
                >
                  {{ record.name }}
                </a-avatar>
              </template>
              <template v-else-if="item.key === 'sex'" #cell="{ record }">
                <a-tag v-if="record.sex === 1" color="blue" size="small">男</a-tag>
                <a-tag v-else color="blue" size="small">女</a-tag>
              </template>
              <template v-else-if="item.key === 'status'" #cell="{ record }">
                <a-tag color="purple" size="small">请假</a-tag>
              </template>
              <template v-else-if="item.key === 'pass'" #cell="{ record }">
                <a-tag v-if="record.pass === 1" color="blue" size="small">通过</a-tag>
                <a-tag v-else-if="record.pass === 2" color="red" size="small">未通过</a-tag>
                <a-tag v-else color="purple" size="small">未审核</a-tag>
              </template>
              <template v-else-if="item.key === 'url'" #cell="{ record }">
                <a-image
                    :preview-props="{
                     actionsLayout: ['rotateRight', 'zoomIn', 'zoomOut'],
               }"
                    :src="record.url"
                    width="100"
                />
              </template>
            </a-table-column>
          </template>
        </a-table>
      </template>
      <template #footer>
        <TableFooter :pagination="pagination"/>
      </template>
    </TableBody>
  </div>
</template>

<script lang="ts">
import {get} from '@/api/http'
import {getExamine,} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn,} from '@/hooks/table'
import FormRender from '@/components/FormRender'
import {FormItem, ModalDialogType} from '@/types/components'
import {Form, Message} from '@arco-design/web-vue'
import {defineComponent, onMounted, ref} from 'vue'
import type {Dayjs} from 'dayjs'
import useUserStore from "@/store/modules/user";
import ModalDialog from "@/components/ModalDialog.vue";

const conditionItems: Array<FormItem> = [
  {
    key: 'time',
    label: '时间',
    type: 'time',
    value: ref<string>(''),
  },
]
export default defineComponent({
  name: 'TableWithSearch',
  components: {
    ModalDialog,
    FormRender,
  },
  setup() {
    const searchForm = ref()
    const pagination = usePagination(doRefresh)
    const {selectedRowKeys, onSelectionChange, showCheckedAll} = useRowSelection()
    const table = useTable()
    const rowKey = useRowKey('id')
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '员工姓名',
        key: 'name',
        dataIndex: 'name',
      },
      {
        title: '头像',
        key: 'avatar',
        dataIndex: 'avatar',
      },
      {
        title: '工号',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '性别',
        key: 'sex',
        dataIndex: 'sex',
      },
      {
        title: '部门',
        key: 'department_id',
        dataIndex: 'department_id',
      },
      {
        title: '开始时间',
        key: 'start_time',
        dataIndex: 'start_time',
      },
      {
        title: '结束时间',
        key: 'end_time',
        dataIndex: 'end_time',
      },
      {
        title: '类型',
        key: 'status',
        dataIndex: 'status',
      },
      {
        title: '图片',
        key: 'url',
        dataIndex: 'url',
      },
      {
        title: '审核状态',
        key: 'pass',
        dataIndex: 'pass',
      },
      {
        title: '说明',
        key: 'reason',
        dataIndex: 'reason',
      },
    ])
    const userStore = useUserStore()
    const actionTitle = ref('请假审核')
    const modalDialogRef = ref<ModalDialogType | null>(null)
    const formItems = [
      {
        label: '姓名',
        key: 'name',
        type: 'input',
        placeholder: '请选择性别',
        value: ref(null),
      },
      {
        label: '工号',
        key: 'uid',
        type: 'input',
        placeholder: '请选择工号',
        value: ref(undefined),
      },
      {
        label: '性别',
        key: 'sex',
        type: 'input',
        placeholder: '请输入性别',
        value: ref(null),
      },
      {
        label: '起止日期',
        key: 'startEndDate',
        type: 'date-range',
        value: ref<Dayjs[]>([]),
        reset: function () {
          this.value.value = []
        },
      },
      {
        label: '请假原因',
        key: 'reason',
        placeholder: '请输入请假原因',
        type: 'textarea',
        value: ref(null),
      },
      {
        label: '是否通过',
        key: 'pass',
        type: 'switch',
        value: ref(false),
        placeholder: '请审核',

      },
    ] as FormItem[]
    const submitLoading = ref(false)
    const department = {
      1: "程序部",
      2: "人事部",
      3: "财务部",
      4: "销售部",
      5: "法务部"
    }

    let formRef = ref<typeof Form>()

    function doRefresh() {
      get({
        url: getExamine,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
      }).then((res) => {
        res.data.forEach((i: any) => {
          const date = new Date(i.start_time);
          const year = date.getFullYear();
          const month = String(date.getMonth() + 1).padStart(2, "0");
          const day = String(date.getDate()).padStart(2, "0");
          const date1 = new Date(i.end_time);
          const year1 = date.getFullYear();
          const month1 = String(date.getMonth() + 1).padStart(2, "0");
          const day1 = String(date.getDate()).padStart(2, "0");
          i.end_time = `${year1}-${month1}-${day1}  ${date1.getHours() >= 10 ? date1.getHours() : '0' + date1.getHours()}:${date1.getMinutes() > 10 ? date1.getMinutes() : '0' + date1.getMinutes()}:${date1.getSeconds() > 10 ? date1.getSeconds() : '0' + date1.getSeconds()}`
          i.start_time = `${year}-${month}-${day}  ${date.getHours() >= 10 ? date.getHours() : '0' + date.getHours()}:${date.getMinutes() > 10 ? date.getMinutes() : '0' + date.getMinutes()}:${date.getSeconds() > 10 ? date.getSeconds() : '0' + date.getSeconds()}`
          return
        })
        table.handleSuccess(res)
        pagination.setTotalSize(res.data.length || 10)
      }).catch(error => {
        Message.error(error.message)
        table.tableLoading.value = false
      })
    }

    function onSearch() {
      let data: any = conditionItems.reduce((pre, cur) => {
        ;(pre as any)[cur.key] = cur.value.value
        return pre
      }, {})
      const tableList = table.dataList.filter(i => {
        if ((new Date(data.time[0]).getTime()) < (new Date(i.start_time).getTime()) && (new Date(i.start_time).getTime()) < (new Date(data.time[1]).getTime())) {
          console.log(i)
          return i
        }
      })
      table.handleSuccess({data: tableList})
      pagination.setTotalSize(tableList.length || 10)
    }


    onMounted(doRefresh)
    return {
      ...table,
      rowKey,
      pagination,
      searchForm,
      tableColumns,
      conditionItems,
      onSearch,
      selectedRowKeys,
      showCheckedAll,
      actionTitle,
      modalDialogRef,
      formItems,
      submitLoading,
      formRef,
      department,
      onSelectionChange,
    }
  },
})
</script>

<style lang="less" scoped>
.form-wrapper {
  font-size: small;
}

.avatar-container {
  position: relative;
  width: 30px;
  height: 30px;
  margin: 0 auto;
  vertical-align: middle;

  .avatar {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }

  .avatar-vip {
    border: 2px solid #cece1e;
  }

  .vip {
    position: absolute;
    top: 0;
    right: -9px;
    width: 15px;
    transform: rotate(60deg);
  }
}

.gender-container {
  .gender-icon {
    width: 20px;
  }
}
</style>
