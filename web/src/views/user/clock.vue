<template>
  <div style="width: 100%; height:100%;"/>
  <TableBody>
    <template #header>
      <TableHeader ref="tableHeaderRef" :show-filter="false" title="考勤列表">
        <template #table-config>
          <a-button size="mini" type="primary" @click="onAddItem"> 打卡</a-button>
        </template>
      </TableHeader>
    </template>
    <template #default>
      <a-table
          :bordered="{ wrapper: true, cell: true }"
          :data="dataList"
          :hoverable="true"
          :loading="tableLoading"
          :pagination="false"
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
              :width="item.width"
          >
            <template v-if="item.key === 'index'" #cell="{ rowIndex }">
              {{ rowIndex + 1 }}
            </template>

            <template v-else-if="item.key === 'status'" #cell="{ record }">
              <a-tag v-if="record.status ===2 " color="blue" size="small"> {{ map[record.status] }}</a-tag>
              <a-tag v-else color="red" size="small"> {{ map[record.status] }}</a-tag>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </template>

    <template #footer>
      <TableFooter ref="tableFooterRef" :pagination="pagination"/>
    </template>
  </TableBody>
</template>

<script lang="ts">
import {clockIn, delete_employer, getClockInLog} from '@/api/url'
import {usePagination, useRowKey, useRowSelection, useTable, useTableColumn, useTableHeight,} from '@/hooks/table'
import {Message, Modal} from '@arco-design/web-vue'
import {defineComponent, getCurrentInstance, onMounted, ref} from 'vue'
import AddButton from "@/components/AddButton.vue";
import useUserStore from "@/store/modules/user";
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
      1: "缺勤",
      2: "打卡",
      3: "迟到",
      4: "加班 ",
      5: "补卡",
      6: "出差",
      7: "请假"
    }
    const tableColumns = useTableColumn([
      table.indexColumn,
      {
        title: '工号',
        key: 'uid',
        dataIndex: 'uid',
      },
      {
        title: '开始时间',
        key: 'start_time',
        dataIndex: 'start_time',
      },
      {
        title: '下班打卡',
        key: 'end_time',
        dataIndex: 'end_time',
      },
      {
        title: '日期',
        key: 'date',
        dataIndex: 'date',
      },
      {
        title: '状态',
        key: 'status',
        dataIndex: 'status',
      },
    ])
    const expandAllFlag = ref(true)
    const userStore = useUserStore()
    const get = useGet()
    let add = true


    function doRefresh() {
      get({
        url: getClockInLog,
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
          i.start_time = GetTime(i.start_time)
          i.end_time = GetTime(i.end_time)
          i.date = Getdate(i.date)
        })
        table.handleSuccess(res)
        pagination.setTotalSize(table.dataList.length)
      }).catch(error => {
        console.log(error)
      })
    }

    function Getdate(item: any) {
      const date = new Date(item);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      return `${year}-${month}-${day} `
    }

    function GetTime(item: any) {
      const date = new Date(item);
      const year = date.getFullYear();
      const month = String(date.getMonth() + 1).padStart(2, "0");
      const day = String(date.getDate()).padStart(2, "0");
      const hours = String(date.getHours()).padStart(2, "0");
      const minutes = String(date.getMinutes()).padStart(2, "0");
      const seconds = String(date.getSeconds()).padStart(2, "0");
      return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
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


    function onAddItem() {
      // 打卡
      get({
        url: clockIn,
        headers: {
          Authorization: "Bearer " + userStore.token
        },
      }).then(res => {
        doRefresh()
        Message.success("打卡成功")
      }).catch(error => {
        Message.error(error.message)
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
      onAddItem,
      onDeleteItem,
      GetTime,
      Getdate,
      map,
      add,
    }
  },
})
</script>
