import './App.css'
import { AppShell, Group, Text, Card, Title, Badge, SimpleGrid } from '@mantine/core'
import { Navbar } from './components/Navbar'

function App() {
  return (
    <AppShell
      navbar={{ width: 250, breakpoint: 'sm' }}
      padding="md"
    >
      <Navbar
        current="overview"
        onNavigate={() => {}}
        onUpload={() => {}}
        monthStatus={{
          label: 'April 2026',
          uploaded: 3,
          expected: 5,
          detail: '2 statements still missing',
        }}
        statementsMissing={2}
        user={{
          name: 'James',
          email: 'james@example.com',
        }}
      />
      <AppShell.Main>
        <Group h={60} px="md" style={{ borderBottom: '1px solid #D4A853' }}>                                                                                                                                     
          <Text fw={400} size="xl">Good morning, James</Text>              
        </Group>
        <SimpleGrid cols={4} mt="md">
          <Card withBorder radius="md" p="md">                                                                                                                                                                       
            <Text size="sm" c="dimmed">Total Spent This Month</Text>                                                                                                                                                 
            <Group justify="space-between" mt="xs">
              <Title order={2}>$1,240</Title>                                                                                                                                                                        
              <Badge color="red">↑ 8%</Badge>     
            </Group>                                                                                                                                                                                                 
          </Card>
          <Card withBorder radius="md" p="md">                                                                                                                                                                       
            <Text size="sm" c="dimmed">Total Spent This Month</Text>                                                                                                                                                 
            <Group justify="space-between" mt="xs">
              <Title order={2}>$1,240</Title>                                                                                                                                                                        
              <Badge color="red">↑ 8%</Badge>     
            </Group>                                                                                                                                                                                                 
          </Card>
          <Card withBorder radius="md" p="md">                                                                                                                                                                       
            <Text size="sm" c="dimmed">Total Spent This Month</Text>                                                                                                                                                 
            <Group justify="space-between" mt="xs">
              <Title order={2}>$1,240</Title>                                                                                                                                                                        
              <Badge color="red">↑ 8%</Badge>     
            </Group>                                                                                                                                                                                                 
          </Card>
          <Card withBorder radius="md" p="md">                                                                                                                                                                       
            <Text size="sm" c="dimmed">Total Spent This Month</Text>                                                                                                                                                 
            <Group justify="space-between" mt="xs">
              <Title order={2}>$1,240</Title>                                                                                                                                                                        
              <Badge color="red">↑ 8%</Badge>     
            </Group>                                                                                                                                                                                                 
          </Card>
        </SimpleGrid>  
      </AppShell.Main>
    </AppShell>
  )
}

export default App
